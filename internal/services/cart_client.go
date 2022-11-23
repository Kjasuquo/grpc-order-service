package services

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"order_svc/proto"
)

func InitCartServiceClient(url string) proto.CartServiceClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect to Cart Client:", err)
	}

	return proto.NewCartServiceClient(cc)
}

func (s *OrderServiceServer) DeleteFromCart(req *proto.CartRequestItem) (*proto.CartResponse, error) {
	return s.CartClient.DeleteCart(context.Background(), req)
}
