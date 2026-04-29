FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o operator ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/operator .

ENTRYPOINT ["./operator"]

