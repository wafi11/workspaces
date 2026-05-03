package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	v1 "github.com/wafi11/workspaces/core/storage-service/gen/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	minioClient *minio.Client
	v1.UnimplementedStorageServiceServer
}

func NewService(minioClient *minio.Client) *Service {
	return &Service{
		minioClient: minioClient,
	}
}

func (s *Service) PostStorage(ctx context.Context, req *v1.PostStorageRequest) (*v1.PostStorageResponse, error) {
	fmt.Printf("request incoming")
	if err := validatePostRequest(req); err != nil {
		return nil, err
	}

	if len(req.FileData) > maxSingleUpload {
		return nil, status.Errorf(codes.InvalidArgument,
			"file terlalu besar (%d bytes), gunakan PostStorageStream untuk file > 4MB", len(req.FileData))
	}

	if err := s.ensureBucket(ctx, req.BucketName); err != nil {
		return nil, err
	}

	contentType := req.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	info, err := s.minioClient.PutObject(ctx, req.BucketName, req.FileName,
		bytes.NewReader(req.FileData),
		int64(len(req.FileData)),
		minio.PutObjectOptions{
			ContentType:  contentType,
			UserMetadata: req.Metadata,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "upload ke minio gagal: %v", err)
	}

	url, err := s.buildURL(ctx, req.BucketName, req.FileName, defaultPresignExpiry)
	if err != nil {
		return nil, err
	}

	return &v1.PostStorageResponse{
		Url:        url,
		BucketName: req.BucketName,
		FileName:   req.FileName,
		Etag:       info.ETag,
		Size:       info.Size,
	}, nil
}

// ---------------------------------------------------------------------------
// PostStorageStream — streaming upload untuk file besar
// ---------------------------------------------------------------------------

func (s *Service) PostStorageStream(stream v1.StorageService_PostStorageStreamServer) error {
	ctx := stream.Context()

	// Chunk pertama harus meta
	first, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "gagal recv chunk pertama: %v", err)
	}

	meta, ok := first.Data.(*v1.PostStorageStreamRequest_Meta)
	if !ok {
		return status.Error(codes.InvalidArgument, "chunk pertama harus berisi metadata")
	}
	m := meta.Meta

	if m.BucketName == "" || m.FileName == "" {
		return status.Error(codes.InvalidArgument, "bucket_name dan file_name wajib diisi di metadata")
	}

	if err := s.ensureBucket(ctx, m.BucketName); err != nil {
		return err
	}

	// Pipe streaming gRPC → MinIO pakai io.Pipe supaya ga perlu buffer semua di RAM
	pr, pw := io.Pipe()

	// Goroutine: terima semua chunk dari gRPC, tulis ke pipe
	var recvErr error
	go func() {
		defer pw.Close()
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				recvErr = err
				pw.CloseWithError(err)
				return
			}
			chunk, ok := msg.Data.(*v1.PostStorageStreamRequest_Chunk)
			if !ok {
				continue
			}
			if _, err := pw.Write(chunk.Chunk); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()

	contentType := m.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// MinIO PutObject baca dari pipe — size -1 = unknown (pakai multipart otomatis)
	info, err := s.minioClient.PutObject(ctx, m.BucketName, m.FileName, pr, m.TotalSize,
		minio.PutObjectOptions{
			ContentType:  contentType,
			UserMetadata: m.Metadata,
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "stream upload ke minio gagal: %v", err)
	}
	if recvErr != nil {
		return status.Errorf(codes.Internal, "error saat terima stream: %v", recvErr)
	}

	url, err := s.buildURL(ctx, m.BucketName, m.FileName, defaultPresignExpiry)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&v1.PostStorageResponse{
		Url:        url,
		BucketName: m.BucketName,
		FileName:   m.FileName,
		Etag:       info.ETag,
		Size:       info.Size,
	})
}

// ---------------------------------------------------------------------------
// GetStorage — presigned URL
// ---------------------------------------------------------------------------

func (s *Service) GetStorage(ctx context.Context, req *v1.GetStorageRequest) (*v1.GetStorageResponse, error) {
	if req.FileName == "" || req.BucketName == "" {
		return nil, status.Error(codes.InvalidArgument, "file_name dan bucket_name wajib diisi")
	}

	expiry := defaultPresignExpiry
	if req.ExpiresIn > 0 {
		expiry = time.Duration(req.ExpiresIn) * time.Second
	}

	// Ambil object info dulu
	stat, err := s.minioClient.StatObject(ctx, req.BucketName, req.FileName, minio.StatObjectOptions{})
	if err != nil {
		minioErr := minio.ToErrorResponse(err)
		if minioErr.Code == "NoSuchKey" {
			return nil, status.Errorf(codes.NotFound, "file %q tidak ditemukan di bucket %q", req.FileName, req.BucketName)
		}
		return nil, status.Errorf(codes.Internal, "stat object gagal: %v", err)
	}

	url, err := s.buildURL(ctx, req.BucketName, req.FileName, expiry)
	if err != nil {
		return nil, err
	}

	return &v1.GetStorageResponse{
		Url:         url,
		FileName:    req.FileName,
		BucketName:  req.BucketName,
		ContentType: stat.ContentType,
		Size:        stat.Size,
		Etag:        stat.ETag,
	}, nil
}

// ---------------------------------------------------------------------------
// DeleteStorage
// ---------------------------------------------------------------------------

func (s *Service) DeleteStorage(ctx context.Context, req *v1.DeleteStorageRequest) (*v1.DeleteStorageResponse, error) {
	if req.FileName == "" || req.BucketName == "" {
		return nil, status.Error(codes.InvalidArgument, "file_name dan bucket_name wajib diisi")
	}

	err := s.minioClient.RemoveObject(ctx, req.BucketName, req.FileName, minio.RemoveObjectOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hapus object gagal: %v", err)
	}

	return &v1.DeleteStorageResponse{
		Success: true,
		Message: fmt.Sprintf("file %q berhasil dihapus dari bucket %q", req.FileName, req.BucketName),
	}, nil
}

// ---------------------------------------------------------------------------
// ListStorage
// ---------------------------------------------------------------------------

func (s *Service) ListStorage(ctx context.Context, req *v1.ListStorageRequest) (*v1.ListStorageResponse, error) {
	if req.BucketName == "" {
		return nil, status.Error(codes.InvalidArgument, "bucket_name wajib diisi")
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 || pageSize > 1000 {
		pageSize = 100
	}

	opts := minio.ListObjectsOptions{
		Prefix:    req.Prefix,
		Recursive: true,
	}

	var objects []*v1.StorageObject
	count := 0

	for obj := range s.minioClient.ListObjects(ctx, req.BucketName, opts) {
		if obj.Err != nil {
			return nil, status.Errorf(codes.Internal, "list objects gagal: %v", obj.Err)
		}
		// Simple pagination: skip sampai page_token
		if req.PageToken != "" && obj.Key <= req.PageToken {
			continue
		}
		objects = append(objects, &v1.StorageObject{
			FileName:    obj.Key,
			BucketName:  req.BucketName,
			Size:        obj.Size,
			ContentType: obj.ContentType,
			Etag:        obj.ETag,
			// LastModified: timestamppb.New(obj.LastModified),
		})
		count++
		if count >= pageSize {
			break
		}
	}

	var nextToken string
	if len(objects) == pageSize {
		nextToken = objects[len(objects)-1].FileName
	}

	return &v1.ListStorageResponse{
		Objects:       objects,
		NextPageToken: nextToken,
		TotalCount:    int32(len(objects)),
	}, nil
}
