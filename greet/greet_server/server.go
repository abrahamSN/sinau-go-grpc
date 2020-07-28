package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/abrahamSN/sinau-go-grpc/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, request *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v \n", request)
	firstName := request.GetGreeting().GetFirstName()

	result := "Hello " + firstName

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(request *greetpb.GreetManyTimeRequest, timesServer greetpb.GreetService_GreetManyTimesServer) error {
	firstName := request.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimeResponse{
			Result: result,
		}

		timesServer.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (*server) LongGreet(greetServer greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked with a streaming request \n")
	result := ""
	for {
		req, err := greetServer.Recv()

		if err == io.EOF {
			// we have finished reading the client stream2
			return greetServer.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Printf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()

		result += "Hello " + firstName + "! \n"
	}
}

func (*server) GreetEveryone(everyoneServer greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked with a streaming request \n")

	for {
		req, err := everyoneServer.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Printf("Error while reading client stream: %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()

		result := "Hello " + firstName + "!\n"

		sendErr := everyoneServer.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if sendErr != nil {
			log.Printf("Error while sending data to client: %v", sendErr)
			return sendErr
		}

	}
}

func main() {
	fmt.Println("Greet Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v ", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
