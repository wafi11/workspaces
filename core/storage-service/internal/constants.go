package internal

import "time"

const (
	defaultPresignExpiry = time.Hour
	maxSingleUpload      = 4 * 1024 * 1024
	defaultChunkSize     = 512 * 1024
)
