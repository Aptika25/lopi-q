FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY proto ./proto
COPY pkg ./pkg
COPY services/auth-service ./services/auth-service

WORKDIR /app/services/auth-service

RUN go mod download

RUN go build -o auth-service ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/auth-service/auth-service .

EXPOSE 50051

CMD ["./auth-service"]
