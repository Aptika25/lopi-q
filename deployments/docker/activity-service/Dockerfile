FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY proto ./proto
COPY pkg ./pkg
COPY services/activity-service ./services/activity-service

WORKDIR /app/services/activity-service

RUN go mod download

RUN go build -o activity-service ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/activity-service/activity-service .

EXPOSE 50053

CMD ["./activity-service"]
