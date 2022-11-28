package main

import (
	"log"
	"net"

	"github.com/daluisgarcia/golang-probuffers-grpc/database"
	serverPkg "github.com/daluisgarcia/golang-probuffers-grpc/server"
	"github.com/daluisgarcia/golang-probuffers-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":5061")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	repo, err := database.NewPostgresRepository("postgres://postgres@localhost:5432/golang-protobuf?sslmode=disable")

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	server := serverPkg.NewTestServer(repo)

	s := grpc.NewServer()

	testpb.RegisterTestServiceServer(s, server)

	reflection.Register(s) // Reflection provides metadata to clients

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
