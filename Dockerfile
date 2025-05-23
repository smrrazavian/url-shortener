FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/api ./cmd/api

FROM alpine:latest

# COPY --from=builder /go/bin/api /app/api
CMD ["./main"]
