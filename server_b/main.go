package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "MockingDemo" // Change to match the module name in go.mod

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type serverB struct {
	pb.UnimplementedServiceBServer
	environment string
}

func (s *serverB) ProcessRequest(ctx context.Context, req *pb.EnhancedRequest) (*pb.ProcessedResponse, error) {
	log.Printf("[%s] Server B received request: %v", s.environment, req)

	// Process the request
	processingTime := time.Now().Format(time.RFC3339Nano)

	// Create response
	response := &pb.ProcessedResponse{
		Result:         fmt.Sprintf("Processed query: %s for user: %s in %s environment", req.Query, req.UserId, s.environment),
		Success:        true,
		ProcessingTime: processingTime,
		ServerBId:      "server-b-instance-001",
	}

	log.Printf("[%s] Server B sending response: %v", s.environment, response)
	return response, nil
}

func main() {
	// Parse command line arguments
	environment := flag.String("env", "stg", "Environment (stg or qa)")
	flag.Parse()

	// Validate environment
	if *environment != "stg" && *environment != "qa" {
		log.Fatalf("Invalid environment: %s. Must be 'stg' or 'qa'", *environment)
	}

	log.Printf("Starting Server B in %s environment", *environment)

	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("[%s] Failed to listen: %v", *environment, err)
	}

	s := grpc.NewServer()
	pb.RegisterServiceBServer(s, &serverB{environment: *environment})
	reflection.Register(s) // Enable reflection for tools like grpcurl

	log.Printf("[%s] Server B listening on %s", *environment, port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", *environment, err)
	}
}
