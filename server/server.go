package main

import (
	"context"
	"fmt"
	"log"
	"net"
	t "time"

	gRPC "github.com/ThomasBavn/Security-Miniproject2/proto"

	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedGetCurrentTimeServer
}

func (s *Server) GetTime(ctx context.Context, in *gRPC.GetTimeRequest) (*gRPC.GetTimeReply, error) {
	fmt.Printf("Received GetTime request\n")
	return &gRPC.GetTimeReply{Reply: t.Now().String()}, nil
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	gRPC.RegisterGetCurrentTimeServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
