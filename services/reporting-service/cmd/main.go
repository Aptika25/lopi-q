package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	repProto "proto/reporting"
	"reporting-service/internal/config"
	"reporting-service/internal/database"
	repGrpc "reporting-service/internal/handler/grpc"
	"reporting-service/internal/repository"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewDB(cfg)
	repo := repository.NewReportingRepository(db)
	handler := repGrpc.NewReportingHandler(repo)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	s := grpc.NewServer()
	repProto.RegisterReportingServiceServer(s, handler)

	fmt.Printf("===================================================\n")
	fmt.Printf("📊 LOPI-Q REPORTING-SERVICE (gRPC) Running on port %s\n", cfg.Port)
	if db != nil {
		fmt.Printf("🐘 Connected to PostgreSQL Database: db_lopiq_reporting\n")
	} else {
		fmt.Printf("📦 Running with In-Memory Repository Fallback\n")
	}
	fmt.Printf("===================================================\n")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
