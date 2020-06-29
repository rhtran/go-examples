package main

import (
	"context"
	"fmt"
	"github.com/rhtran/go-examples/grpc-course/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {}

func(*server) Greet(ctx context.Context, request *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", request)
	firstName := request.Greeting.FirstName
	result := "Hello " + firstName

	response := &greetpb.GreetResponse{
		Result: result,
	}

	return response, nil
}

func (*server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", request)
	firstName := request.Greeting.FirstName
	repeats := request.Repeats

	for i := 0; int64(i) < repeats; i++ {
		response := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + firstName + " time: " + strconv.Itoa(i),
		}
		stream.Send(response)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with streaming request\n")

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		firstName := request.Greeting.FirstName
		result := "Hello " + firstName + "!"
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
		}
	}

}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	result := ""
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			// finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := request.Greeting.FirstName
		result += "Hello " + firstName + "!\n"
	}

	return nil
}

func main() {
	fmt.Println("Hello protobuf")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}


