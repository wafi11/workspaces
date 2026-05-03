# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod files first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates (buat HTTPS calls, OAuth, dll)
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/server .

EXPOSE 5051

CMD ["./server"]
