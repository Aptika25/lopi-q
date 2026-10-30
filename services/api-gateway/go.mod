module api-gateway

go 1.23

require (
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.32.0
	pkg v0.0.0
	proto v0.0.0
)

require (
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/pquerna/otp v1.4.0
)

replace pkg => ../../pkg

replace proto => ../../proto
