package internal

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	v1 "github.com/wafi11/workspaces/core/storage-service/gen/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// ensureBucket buat bucket kalau belum ada
func (s *Service) ensureBucket(ctx context.Context, bucket string) error {
	exists, err := s.minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return status.Errorf(codes.Internal, "cek bucket gagal: %v", err)
	}
	if !exists {
		if err := s.minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return status.Errorf(codes.Internal, "buat bucket gagal: %v", err)
		}
	}
	return nil
}

// buildURL generate presigned URL; kalau ada PublicEndpoint ganti host-nya
func (s *Service) buildURL(ctx context.Context, bucket, file string, expiry time.Duration) (string, error) {
	u, err := s.minioClient.PresignedGetObject(ctx, bucket, file, expiry, nil)
	if err != nil {
		return "", status.Errorf(codes.Internal, "presign URL gagal: %v", err)
	}

	return u.String(), nil
}

func validatePostRequest(req *v1.PostStorageRequest) error {
	if req.BucketName == "" {
		return status.Error(codes.InvalidArgument, "bucket_name wajib diisi")
	}
	if req.FileName == "" {
		return status.Error(codes.InvalidArgument, "file_name wajib diisi")
	}
	if len(req.FileData) == 0 {
		return status.Error(codes.InvalidArgument, "file_data tidak boleh kosong")
	}
	return nil
}
