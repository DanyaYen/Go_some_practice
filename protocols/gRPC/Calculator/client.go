package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"calculator/calculator" 
)

const (
	address = "localhost:50051"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: client <num1> <num2>")
		return
	}

	num1, err1 := strconv.Atoi(os.Args[1])
	num2, err2 := strconv.Atoi(os.Args[2])
	if err1 != nil || err2 != nil {
		log.Fatalf("Invalid arguments: %v, %v", err1, err2)
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := calculator.NewCalculatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sumResult, err := c.Add(ctx, &calculator.CalcRequest{FirstNumber: int32(num1), SecondNumber: int32(num2)})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Printf("Sum: %d", sumResult.GetResult())

	subtractResult, err := c.Subtract(ctx, &calculator.CalcRequest{FirstNumber: int32(num1), SecondNumber: int32(num2)})
	if err != nil {
		log.Fatalf("could not subtract: %v", err)
	}
	log.Printf("Subtract: %d", subtractResult.GetResult())
}