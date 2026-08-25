module api-gateway

go 1.26.1

require (
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.32.0
	pkg v0.0.0
	proto v0.0.0
)

require github.com/golang-jwt/jwt/v5 v5.2.1 // indirect

replace pkg => ../../pkg

replace proto => ../../proto
