FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY proto ./proto
COPY pkg ./pkg
COPY services/reporting-service ./services/reporting-service

WORKDIR /app/services/reporting-service

RUN go mod download

RUN go build -o reporting-service ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/reporting-service/reporting-service .

EXPOSE 50054

CMD ["./reporting-service"]
