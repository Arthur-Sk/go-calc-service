package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
	//calcAverage(c)
	//findMax(c)
	doErrorUnary(c)
	doSquareRootUnary(c)
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
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("The sum of %d and %d is: %v", first, second, resp.Sum)
}

func doErrorUnary(c calcpb.CalcServiceClient) {
	rand.Seed(time.Now().UnixNano())
	num := -rand.Intn(100)

	doSquareRootCall(int64(num), c)
}

func doSquareRootUnary(c calcpb.CalcServiceClient) {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(100)

	doSquareRootCall(int64(num), c)
}

func doSquareRootCall(num int64, c calcpb.CalcServiceClient) {
	req := &calcpb.SquareRootRequest{
		Num: num,
	}

	resp, err := c.SquareRoot(context.Background(), req)
	if err == nil {
		log.Printf("Result of square root of %v is: %v", num, resp.NumberRoot)

		return
	}

	respErr, ok := status.FromError(err)
	if false == ok {
		log.Fatalf("Fatal error while calling SquareRoot RPC: %v", err)
	}

	// actual error from gRPC (user error)
	fmt.Println("Error from server: " + respErr.Message())
	fmt.Println(respErr.Code())
	if respErr.Code() == codes.InvalidArgument {
		fmt.Println("We probably sent a genative number")
	}
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

func findMax(c calcpb.CalcServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while calling FindMaximum: %v", err)
	}

	iterations := rand.Intn(25)

	requests := []*calcpb.FindMaximumRequest{}
	for i := 0; i < iterations; i++ {
		requests = append(requests, &calcpb.FindMaximumRequest{
			Num: int64(rand.Intn(100000)),
		})
	}

	waitC := make(chan struct{})

	// Block until everything is done
	defer func() {
		<-waitC
	}()

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(500 * time.Millisecond)
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
				log.Fatalf("Error while receiving: %v", err)
			}
			fmt.Printf("received: %v\n", res)
		}

		close(waitC)
	}()
}
