package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/abrahamSN/sinau-go-grpc/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {}

func (*server) Sum(ctx context.Context, request *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Received Sum RPC: %v \n", request)

	firstNumber := request.GetFirstNumber()
	secondNumber := request.GetSecondNumber()

	sum := firstNumber + secondNumber

	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}

	return res, nil
}

func (*server) PrimeNumber(request *calculatorpb.PrimeNumberRequest, numberServer calculatorpb.CalculatorService_PrimeNumberServer) error {
	fmt.Printf("Received PrimeNumber RPC: %v \n", request)

	number := request.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number % divisor == 0 {
			numberServer.Send(&calculatorpb.PrimeNumberResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has incrised to: %v \n", divisor)
		}
	}
	return nil
}

func (*server) ComputeAverage(averageServer calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Received Compute Average RPC: \n")

	sum := int32(0)
	count := 0

	for {
		req, err := averageServer.Recv()

		if err == io.EOF {
			// we have finished reading the client stream2
			averageData := float64(sum) / float64(count)
			return averageServer.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: averageData,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		sum += req.GetNumber()
		count++

		fmt.Printf("ini sum: %v \n", sum)
	}
}

func (*server) FindMaximum(maximumServer calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("Received FindMaximum RPC: \n")

	maximum := int32(0)

	for {
		req, err := maximumServer.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		number := req.GetNumber()

		if number > maximum {
			maximum = number

			sendErr := maximumServer.Send(&calculatorpb.FindMaximumResponse{
				Maximum: maximum,
			})

			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
			}
		}
	}
}

func main() {
	fmt.Printf("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v ", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
