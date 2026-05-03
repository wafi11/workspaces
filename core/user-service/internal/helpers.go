package internal

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func detectContentType(data []byte) (mimeType string, ext string, err error) {
	if len(data) < 4 {
		return "", "", status.Error(codes.InvalidArgument, "file terlalu kecil")
	}

	switch {
	// JPEG: FF D8 FF
	case data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF:
		return "image/jpeg", ".jpg", nil

	// PNG: 89 50 4E 47 0D 0A 1A 0A
	case data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47:
		return "image/png", ".png", nil

	// WebP: RIFF....WEBP
	case len(data) >= 12 &&
		data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50:
		return "image/webp", ".webp", nil

	default:
		return "", "", status.Error(codes.InvalidArgument,
			"format tidak didukung, hanya jpg/png/webp yang diizinkan")
	}
}
