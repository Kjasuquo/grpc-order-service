package main

import (
	"log"
	"order_svc/config"
	"order_svc/internal/rabbitMQ"
	myService "order_svc/internal/services"
	// "log"
	// "github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	configs := config.ReadConfigs(".")
	// start rabbitMQ
	rabbitMQ.ConnectRabbitMq()

	go func() {
		rabbitMQ.DemoReadFromOrderQueue(configs.OrderCreationChannel)
	}()

	myService.Start()
}
