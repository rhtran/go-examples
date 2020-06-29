package main

import (
	"context"
	"fmt"
	"github.com/rhtran/go-examples/grpc-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Calculator client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	doUnary(c)
	doServeStream(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a Sum Unary RPC...")
	request := &calculatorpb.SumRequest{
			FirstNumber: 23,
			SecondNumber:  59,
	}
	response, err := c.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v\n", response.Sum)
}

func doServeStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a PrimeComposition stream RPC...")
	request := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 84793,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), request)
	if err != nil {
		log.Fatalf("error while reading : %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeFactorDecomposition: %v", response.PrimeFactor)
	}
}