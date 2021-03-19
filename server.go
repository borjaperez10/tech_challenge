package main

import (
	"log"
	"net"

	"github.com/borjaperez10/tech_challenge/communication"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	s := communication.Server{}

	grpcServer := grpc.NewServer()

	communication.RegisterCommunicationServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve GRPC server over port 9000: %v", err)
	}

}