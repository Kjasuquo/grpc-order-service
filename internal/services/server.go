package services

import (
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"

	"order_svc/config"
	"order_svc/internal/ports"
	"order_svc/internal/postgres"

	"order_svc/proto"
	//"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc"
)

type OrderServiceServer struct {
	proto.UnimplementedOrderServer
	Order      ports.Order
	CartClient proto.CartServiceClient
}

func Start() {
	//var configStruct config.Config?
	configs := config.ReadConfigs(".")

	PORT := fmt.Sprintf(":%s", configs.GrpcPort)
	if PORT == ":" {
		PORT += "8080"
	}

	db := &postgres.PostgresDB{}
	db.Init()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%v", PORT))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// grpc server
	grpcServer := grpc.NewServer()

	healthpb.RegisterHealthServer(grpcServer, health.NewServer())

	env, err := config.VaultSecrets(viper.GetString("VAULT_ADDR"), viper.GetString("VAULT_AUTH_TOKEN"), viper.GetString("VAULT_SECRET_PATH"))
	if err != nil {
		log.Println("Could not read secrets from vault", err)
		return
	}

	consulClient, err := config.NewConsulClient(env.ConsulAddress)
	if err != nil {
		log.Println("couldn't connect to consul", err)
	}
	fmt.Println("consul address: ", env.ConsulAddress)

	go func() {
		consulClient.ServiceRegistryWithConsul("order-grpc", "order", PORT, []string{"GRPC", "backend"})
		consulClient.ServiceRegistryWithConsul("order-http", "order", ":8205", []string{"HTTP", "envoy"})
	}()

	serviceAddress, err := consulClient.GetConsulService("cart", "GRPC")
	if err != nil {
		log.Println("couldn't get search service address", err)
	}
	log.Println("cart-service address: ", serviceAddress)

	CartClient := InitCartServiceClient(serviceAddress)

	OrderServiceServer := &OrderServiceServer{
		Order:      db,
		CartClient: CartClient,
	}

	proto.RegisterOrderServer(grpcServer, OrderServiceServer)
	log.Printf("Order server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("failed to serve: %v", err)
	}
}
