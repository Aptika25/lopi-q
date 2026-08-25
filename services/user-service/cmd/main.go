package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	userProto "proto/user"
	"user-service/internal/config"
	"user-service/internal/database"
	userGrpc "user-service/internal/handler/grpc"
	"user-service/internal/repository"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewDB(cfg)
	authDB := database.NewAuthDB(cfg)
	repo := repository.NewUserRepository(cfg.DBPath, db, authDB)
	handler := userGrpc.NewUserHandler(repo)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	s := grpc.NewServer()
	userProto.RegisterUserServiceServer(s, handler)

	fmt.Printf("===================================================\n")
	fmt.Printf("👤 LOPI-Q USER-SERVICE (gRPC) Running on port %s\n", cfg.Port)
	if db != nil {
		fmt.Printf("🐘 Connected to PostgreSQL Database: db_lopiq_user\n")
	}
	fmt.Printf("===================================================\n")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
