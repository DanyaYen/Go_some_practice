package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"calculator/calculator"
)

type server struct {
	calculator.UnimplementedCalculatorServer
}

func (s *server) Add(ctx context.Context, req *calculator.CalcRequest) (*calculator.CalcResponse, error) {
	a := req.GetFirstNumber()
	b := req.GetSecondNumber()
	result := a + b
	return &calculator.CalcResponse{Result: result}, nil
}

func (s *server) Subtract(ctx context.Context, req *calculator.CalcRequest) (*calculator.CalcResponse, error) {
	a := req.GetFirstNumber()
	b := req.GetSecondNumber()
	result := a - b
	return &calculator.CalcResponse{Result: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	calculator.RegisterCalculatorServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}