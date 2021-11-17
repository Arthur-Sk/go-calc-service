package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-service/calculator/calcpb"
	"log"
	"net"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calcpb.CalcActionReq) (*calcpb.CalcActionResp, error) {
	fmt.Printf("Summing %v\n", req.Terms)

	var result int32
	for _, v := range req.Terms {
		result += v
	}

	resp := calcpb.CalcActionResp{
		Sum: int64(result),
	}

	fmt.Printf("The sum is: %d\n", resp.Sum)

	return &resp, nil
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
