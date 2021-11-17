package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-service/calculator/calcpb"
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

	doUnary(c)
}

func doUnary(c calcpb.CalcServiceClient) {
	rand.Seed(time.Now().UnixNano())
	first := rand.Intn(100)
	second := rand.Intn(100)

	req := &calcpb.CalcActionReq{
		Terms: []int32{int32(first), int32(second)},
	}

	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("The sum of %d and %d is: %v", first, second, resp.Sum)
}
