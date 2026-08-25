FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY proto ./proto
COPY pkg ./pkg
COPY services/user-service ./services/user-service

WORKDIR /app/services/user-service

RUN go mod download

RUN go build -o user-service ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/user-service/user-service .

EXPOSE 50052

CMD ["./user-service"]
