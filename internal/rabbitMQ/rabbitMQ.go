package rabbitMQ

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"order_svc/config"
	"order_svc/internal/models"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Url      string
}

// type RabbitMQ struct {
// 	Connection *amqp.Connection
// }

var Connection *amqp.Connection

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func ConnectRabbitMq() {
	configs := config.ReadConfigs(".")

	config := &Config{
		Host:     configs.MqHost,
		Password: configs.MqPassword,
		User:     configs.MqUser,
		Port:     configs.MqPort,
		Url:      configs.RabbitMQUrl,
	}
	var rabbitMQURL string
	if config.Url == "" {
		rabbitMQURL = fmt.Sprintf("amqp://%s:%s@%s:%s/", config.User, config.Password, config.Host, config.Port)
	} else {
		rabbitMQURL = config.Url
	}

	fmt.Println("Connecting to RabbitMQ ...")
	fmt.Println("rabbitMQURL:", rabbitMQURL)

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(rabbitMQURL)
	// connectRabbitMQ, err := amqp.Dial("amqp://guest:guest@localhost:5672")

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to RabbitMQ ...")
	// defer connectRabbitMQ.Close()

	Connection = connectRabbitMQ
}

func PublishToOrderNotificationQueue(orderDetails *models.NewOrder) error {

	// fmt.Println("orderDetails: ", orderDetails)
	ch, err := Connection.Channel()
	failOnError(err, "Failed to open a channel at Publish")
	defer ch.Close()

	// err = ch.ExchangeDeclare(
	// 	"order_notification", // name
	// 	"fanout",             // type
	// 	true,                 // durable
	// 	false,                // auto-deleted
	// 	false,                // internal
	// 	false,                // no-wait
	// 	nil,                  // arguments
	// )

	q, err := ch.QueueDeclare(
		"order_notification", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {

		log.Panic(err, "could not declare queue")
	}

	body, err := json.Marshal(orderDetails)
	if err != nil {
		log.Printf("err :%v", err)
		return err
	}

	// err = ch.Publish(
	// 	"order_notification", // exchange
	// 	"",                   // routing key
	// 	false,                // mandatory
	// 	false,                // immediate
	// 	amqp.Publishing{
	// 		ContentType: "json",
	// 		Body:        body,
	// 	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "json",
			Body:        body,
		},
	)

	log.Printf("Sent %s\n", body)

	if err != nil {
		log.Fatalf("err: %v", err)
		return err
	}

	log.Println("published to order_notification queue")
	return nil
}

// DemoReadFromOrderQueue(channelName string) ([]models.NewOrder, error) {
func DemoReadFromOrderQueue(channelName string) {
	// var orderList []models.NewOrder
	// ConnectRabbitMq()
	ch, err := Connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		channelName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		channelName, // exchange
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	if err != nil {

		failOnError(err, "Failed to register a consumer")
	}

	// orderData := &models.NewOrder{}
	var orderData = &models.NewOrder{}
	var forever chan struct{}

	go func() {
		for d := range msgs {
			// log.Printf("our consumer received a message from: %v", os.Getpid())
			err := json.Unmarshal((d.Body), &orderData)
			if err != nil {
				log.Fatalf("error%v", err)
			}
			//orderCreation.OrderCreator(orderData)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
