package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	pb "MockingDemo" // Updated import path since proto file is not in a proto folder

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const (
	serverBAddress = "localhost:50051"
	serverAVersion = "v1.0.0"
	bufSize        = 1024 * 1024
)

// Global variable for the client
var (
	clientB     pb.ServiceBClient
	environment *string
)

// MockServiceBServer implements the ServiceB interface for testing
type MockServiceBServer struct {
	pb.UnimplementedServiceBServer
}

// ProcessRequest implements the ProcessRequest method for the mock server
func (s *MockServiceBServer) ProcessRequest(ctx context.Context, req *pb.EnhancedRequest) (*pb.ProcessedResponse, error) {
	log.Printf("[MOCK] Received request: %v", req)

	// Simulate processing delay
	time.Sleep(100 * time.Millisecond)

	// Return mock response
	return &pb.ProcessedResponse{
		Result:         fmt.Sprintf("Mock response for query: %s", req.Query),
		Success:        true,
		ProcessingTime: time.Now().Format(time.RFC3339Nano),
		ServerBId:      "mock-server-b-instance",
	}, nil
}

// setupMockServer creates a mock gRPC server for testing
func setupMockServer() (pb.ServiceBClient, func()) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterServiceBServer(s, &MockServiceBServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve mock server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("Failed to dial mock server: %v", err)
	}

	closer := func() {
		lis.Close()
		s.Stop()
		conn.Close()
	}

	return pb.NewServiceBClient(conn), closer
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse user request
	var userReq pb.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[%s] Server A received request: %v", *environment, userReq)

	// Enhance the request with additional fields
	requestID := uuid.New().String()
	enhancedReq := &pb.EnhancedRequest{
		UserId:         userReq.UserId,
		Query:          userReq.Query,
		Timestamp:      time.Now().Format(time.RFC3339),
		RequestId:      requestID,
		ServerAVersion: serverAVersion,
	}

	// Call Server B
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Printf("[%s] Server A sending request to Server B: %v", *environment, enhancedReq)
	resp, err := clientB.ProcessRequest(ctx, enhancedReq)
	if err != nil {
		log.Printf("[%s] Error calling Server B: %v", *environment, err)
		http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Server A received response from Server B: %v", *environment, resp)

	// Transform response for the user
	userResp := &pb.UserResponse{
		Result:    resp.Result,
		Success:   resp.Success,
		RequestId: requestID,
	}

	// Add environment info to response headers
	w.Header().Set("X-Environment", *environment)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(userResp); err != nil {
		log.Printf("[%s] Error encoding response: %v", *environment, err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Server A sent response to user: %v", *environment, userResp)
}

func main() {
	// Parse command line arguments
	environment = flag.String("env", "stg", "Environment (stg or qa)")
	flag.Parse()

	// Validate environment
	if *environment != "stg" && *environment != "qa" {
		log.Fatalf("Invalid environment: %s. Must be 'stg' or 'qa'", *environment)
	}

	log.Printf("Starting Server A in %s environment", *environment)

	var cleanup func()

	// Set up connection based on environment
	if *environment == "qa" {
		// Use mock server for QA environment
		log.Println("Using mock Server B for QA environment")
		clientB, cleanup = setupMockServer()
		defer cleanup()
	} else {
		// Use real server for STG environment
		log.Println("Connecting to real Server B for STG environment")
		conn, err := grpc.Dial(serverBAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to Server B: %v", err)
		}
		defer conn.Close()
		clientB = pb.NewServiceBClient(conn)
	}

	// Create HTTP handler for REST API
	http.HandleFunc("/api/process", processHandler)

	// Start HTTP server
	log.Printf("[%s] Server A listening on :8080", *environment)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[%s] Failed to start server: %v", *environment, err)
	}
}
