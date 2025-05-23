package main

import (
	"context"
	"log"
	"net"
	"google.golang.org/grpc"
	"gRPC/greet"
)

type server struct {
	greet.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
    name := req.GetName()
    message := "Hello, " + name + "!"
    return &greet.GreetResponse{Result: message}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051") 
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    greet.RegisterGreetServiceServer(s, &server{}) 
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}