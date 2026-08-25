FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY proto ./proto
COPY pkg ./pkg
COPY services/api-gateway ./services/api-gateway

WORKDIR /app/services/api-gateway

RUN go mod download

RUN go build -o api-gateway ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/api-gateway/api-gateway .

EXPOSE 5000

CMD ["./api-gateway"]
