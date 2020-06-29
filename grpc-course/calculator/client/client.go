package main

import (
	"context"
	"fmt"
	"github.com/rhtran/go-examples/grpc-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
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
	doClientStream(c)
	doBidirectionStream(c)
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

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a ComputeAverge client stream RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 8, 54,23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The average is: %v\n", response.Average)
}

func doBidirectionStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a FindMaximum bidirection stream RPC...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)
	}

	waitc := make(chan struct{})

	// send go routine
	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Printf("Sending number: %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive go routine
	go func() {
		for  {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while reading server stream: %v", err)
				break
			}

			max := response.Max
			fmt.Printf("Received a new max of ...: %v\n", max)
		}
		close(waitc)
	}()

	// block
	<- waitc
}