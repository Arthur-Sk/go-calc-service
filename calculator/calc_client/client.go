package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-service/calculator/calcpb"
	"io"
	"log"
	"math/rand"
	"time"
)

func main() {
	cc, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not coonect: %v", err)
	}

	defer cc.Close()

	c := calcpb.NewCalcServiceClient(cc)

	//doUnary(c)
	//decompose(c)
	calcAverage(c)
}

func decompose(c calcpb.CalcServiceClient) {
	var target int64 = 120

	req := &calcpb.PrimeNumberDecomposeReq{
		Target: target,
	}

	resStream, err := c.DecomposeToPrime(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Received prime number factor: %d of number %d from server", msg.GetComponent(), target)
	}

	log.Println("End of stream")
}

func doUnary(c calcpb.CalcServiceClient) {
	rand.Seed(time.Now().UnixNano())
	first := rand.Intn(100)
	second := rand.Intn(100)

	req := &calcpb.SumActionReq{
		Terms: []int32{int32(first), int32(second)},
	}

	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("The sum of %d and %d is: %v", first, second, resp.Sum)
}

func calcAverage(c calcpb.CalcServiceClient) {
	rand.Seed(time.Now().UnixNano())
	iterations := rand.Intn(10)

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v", err)
	}

	requests := []*calcpb.ComputeAverageRequest{}
	for i := 0; i < iterations; i++ {
		requests = append(requests, &calcpb.ComputeAverageRequest{
			Member: int64(rand.Intn(100)),
		})
	}

	fmt.Printf("Requests: %v", requests)

	// Send messages
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server: %v", err)
	}

	fmt.Printf("Average received: %f", resp.Average)
}
