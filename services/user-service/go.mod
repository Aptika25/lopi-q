module user-service

go 1.26.1

require (
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.26.0
	google.golang.org/grpc v1.64.0
	proto v0.0.0
)

require (
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace pkg => ../../pkg

replace proto => ../../proto
