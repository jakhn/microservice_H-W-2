package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "app/calculatorpb"
)

func main() {

	conn, err := grpc.Dial("localhost:9002", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewCalculatorServiceClient(conn)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDirectionalStream(c)
}

func doUnary(c pb.CalculatorServiceClient) {

	fmt.Println("Starting to find Square Root Unary RPC...")

	req := pb.SquareRootRequest{
		FirstNumber: 10,
	}

	res, err := c.SquareRoot(context.Background(), &req)
	if err != nil {
		log.Println("error while calling Square Root RPC:", err)
	}

	fmt.Println("Reponse from Square Root:", res.GetSqrtResult())
}

func doServerStreaming(c pb.CalculatorServiceClient) {

	fmt.Println("Starting to check the number to the perfection with PerfectNumber Server Streaming RPC...")

	stream, err := c.PerfectNumber(context.Background(), &pb.PerfectNumberRequest{
		Number: 28,
	})
	if err != nil {
		log.Println("error while calling PerfectNumber RPC:", err)
	}

	for {

		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("error while calling PerfectNumber Stream Recv RPC:", err)
		}

		fmt.Println(resp)
	}
}

func doClientStreaming(c pb.CalculatorServiceClient) {

	stream, err := c.TotalNumber(context.Background())
	if err != nil {
		log.Println("error while calling TotalNumber RPC:", err)
	}
	
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, num := range nums {
		err = stream.Send(&pb.TotalNumberRequest{
			Number: int64(num),
		})
		
		if err != nil {
			log.Println("error while calling TotalNumber Client Stream Send RPC:", err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("error while calling TotalNumber Client Stream Send RPC:", err)
		return
	}

	fmt.Println(res)
	
}

func doBiDirectionalStream(c pb.CalculatorServiceClient) {

	stream, err := c.FindMinimum(context.Background())
	if err != nil {
		log.Println("error while calling FindMinimum Send RPC:", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int32{1, 4,-1, 6,-2, 2, 23, -5, 13, 7, 8, 53, -8, -9}

		for _, num := range numbers {
			stream.Send(&pb.FindMinimumRequest{Number: num})
			time.Sleep(time.Second * 1)
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
				log.Println("error while calling FindMinimum Recv RPC:", err)
				return
			}

			fmt.Println(res)
		}

		close(waitc)
	}()

	<-waitc
}
