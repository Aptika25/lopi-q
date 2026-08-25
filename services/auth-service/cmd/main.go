package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"auth-service/internal/config"
	"auth-service/internal/database"
	authGrpc "auth-service/internal/handler/grpc"
	"auth-service/internal/repository"
	authProto "proto/auth"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewDB(cfg)
	userDB := database.NewUserDB(cfg)
	repo := repository.NewUserRepository(cfg.DBPath, db, userDB)
	handler := authGrpc.NewAuthHandler(repo)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	s := grpc.NewServer()
	authProto.RegisterAuthServiceServer(s, handler)

	fmt.Printf("===================================================\n")
	fmt.Printf("🔐 LOPI-Q AUTH-SERVICE (gRPC) Running on port %s\n", cfg.Port)
	if db != nil {
		fmt.Printf("🐘 Connected to PostgreSQL Database: db_lopiq_auth\n")
	}
	fmt.Printf("===================================================\n")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
