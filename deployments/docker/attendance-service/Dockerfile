FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY proto ./proto
COPY pkg ./pkg
COPY services/attendance-service ./services/attendance-service

WORKDIR /app/services/attendance-service

RUN go mod download

RUN go build -o attendance-service ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/attendance-service/attendance-service .

EXPOSE 50055

CMD ["./attendance-service"]
