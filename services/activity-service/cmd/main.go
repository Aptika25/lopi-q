package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	actProto "proto/activity"
	"activity-service/internal/config"
	"activity-service/internal/database"
	actGrpc "activity-service/internal/handler/grpc"
	"activity-service/internal/repository"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewDB(cfg)
	repo := repository.NewActivityRepository(db)
	handler := actGrpc.NewActivityHandler(repo)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	s := grpc.NewServer()
	actProto.RegisterActivityServiceServer(s, handler)

	fmt.Printf("===================================================\n")
	fmt.Printf("📝 LOPI-Q ACTIVITY-SERVICE (gRPC) Running on port %s\n", cfg.Port)
	if db != nil {
		fmt.Printf("🐘 Connected to PostgreSQL Database: db_lopiq_activity\n")
	} else {
		fmt.Printf("📦 Running with In-Memory Repository Fallback\n")
	}
	fmt.Printf("===================================================\n")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
