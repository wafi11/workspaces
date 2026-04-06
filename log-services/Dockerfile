FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk --no-cache add git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o log-service ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/log-service .

EXPOSE 8080

CMD ["./log-service"]