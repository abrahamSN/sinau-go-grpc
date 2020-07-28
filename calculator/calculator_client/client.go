package main

import (
	"context"
	"fmt"
	"github.com/abrahamSN/sinau-go-grpc/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)
	doServerStreaming(c)
	doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting doUnary RPC...")

	req := &calculatorpb.SumRequest{
		FirstNumber: 1,
		SecondNumber: 30,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}

	log.Printf("Response from Sum: %v", res)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting doServerStreaming RPC...")

	req := &calculatorpb.PrimeNumberRequest{
		Number: 200,
	}

	stream, err := c.PrimeNumber(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Prime RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happend: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting doClientStreaming RPC...")

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{1,3,5,6,7,9}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v  \n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, errs := stream.CloseAndRecv()

	if errs != nil {
		log.Fatalf("Error while recieving response: %v \n", errs)
	}

	fmt.Printf("The average is: %v", res.GetAverage())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting doBiDiStreaming RPC...")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int32{1,4,71,8,34,123}

		for _, number := range numbers {
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
			}

			maximum := res.GetMaximum()
			fmt.Printf("Received a new maximum of ...: %v \n", maximum)
		}
		close(waitc)
	}()

	<-waitc

}