package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-service/calculator/calcpb"
	"io"
	"log"
	"net"
	"time"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calcpb.SumActionReq) (*calcpb.SumActionResp, error) {
	fmt.Printf("Summing %v\n", req.Terms)

	var result int32
	for _, v := range req.Terms {
		result += v
	}

	resp := calcpb.SumActionResp{
		Sum: int64(result),
	}

	fmt.Printf("The sum is: %d\n", resp.Sum)

	return &resp, nil
}

func (s *server) DecomposeToPrime(req *calcpb.PrimeNumberDecomposeReq, stream calcpb.CalcService_DecomposeToPrimeServer) error {
	ch := make(chan int64)

	go s.decomposeToPrime(req.Target, ch)

	for component := range ch {
		response := &calcpb.PrimeNumberDecomposeResp{
			Component: component,
		}

		stream.Send(response)
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (*server) decomposeToPrime(target int64, ch chan int64) {
	var k int64 = 2
	for target > 1 {
		if target%k == 0 { // if k evenly divides into target
			ch <- k
			target = target / k // divide target by k so that we have the rest of the number left.
		} else {
			k = k + 1
		}
	}

	close(ch)
}

func (*server) ComputeAverage(stream calcpb.CalcService_ComputeAverageServer) error {
	var numbers []int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// The client stream has been finished
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Average: calcAverage(numbers),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		numbers = append(numbers, req.Member)
	}
}

func calcAverage(numbers []int64) float64 {
	var total int64 = 0
	for _, number := range numbers {
		total += number
	}

	return float64(total) / float64(len(numbers))
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Cannot listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalcServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
