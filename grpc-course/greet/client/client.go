package main

import (
	"context"
	"fmt"
	"github.com/rhtran/go-examples/grpc-course/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello, I am protobuf client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doUnary(c)
	doServerStream(c)
	doClientStream(c)
	doBidirectionStream(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Unary RPC...")
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ethan",
			LastName:  "Tran",
		},
	}
	response, err := c.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", response.Result)
}

func doServerStream(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Server Stream RPC...")
	request := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ethan",
			LastName:  "Tran",
		},
		Repeats: 15,
	}
	responseStream, err := c.GreetManyTimes(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	for {
		message, err := responseStream.Recv()
		if err == io.EOF {
			// end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", message.Result)
	}
}

func doClientStream(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest {
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ryan",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ethan",
			},
		},		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anh",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
	}

	for _, request := range requests {
		fmt.Printf("sending request: %v\n", request)
		stream.Send(request)
		time.Sleep(100 * time.Millisecond)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receive response LongGreet: %v", err)
	}
	log.Printf("Response from LongGreet: \n%v\n", response.Result)
}

func doBidirectionStream(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Bidirectional Streaming RPC...")

	// create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest {
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ryan",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ethan",
			},
		},		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anh",
			},
		},
	}

	waitc := make(chan struct{})
	// send a bunch of messages to the server (go routine)
	go func() {
		for _, request := range requests {
			fmt.Printf("Sending message: %v\n", request)
			stream.Send(request)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive a bunch of message from the server (go routine)
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", response)
		}
		close(waitc)
	}()

	// block until everything is done
	<- waitc
}