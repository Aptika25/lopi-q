package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	attProto "proto/attendance"
	"attendance-service/internal/config"
	"attendance-service/internal/database"
	attGrpc "attendance-service/internal/handler/grpc"
	"attendance-service/internal/repository"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewDB(cfg)
	repo := repository.NewAttendanceRepository(db)
	handler := attGrpc.NewAttendanceHandler(repo)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	s := grpc.NewServer()
	attProto.RegisterAttendanceServiceServer(s, handler)

	fmt.Printf("===================================================\n")
	fmt.Printf("⏱️ LOPI-Q ATTENDANCE-SERVICE (gRPC) Running on port %s\n", cfg.Port)
	if db != nil {
		fmt.Printf("🐘 Connected to PostgreSQL Database: db_lopiq_attendance\n")
	} else {
		fmt.Printf("📦 Running with In-Memory Repository Fallback\n")
	}
	fmt.Printf("===================================================\n")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
