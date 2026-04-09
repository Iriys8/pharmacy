package main

import (
	"fmt"
	"log"
	"os"

	shared_models "pharmacy-api/shared/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func setupBroker() *shared_models.Broker {
	connString := fmt.Sprintf("amqp://%s:%s@%s:5672/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
	)
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatal("failed to connect to broker:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatal("failed to open channel:", err)
	}

	err = ch.ExchangeDeclare("public_exchange", "direct", true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		log.Fatal("failed to declare exchange:", err)
	}

	queues := []struct {
		name       string
		routingKey string
	}{
		{"public_goods_queue", "goods"},
		{"public_schedule_queue", "schedule"},
		{"public_orders_queue", "orders"},
		{"public_announces_queue", "announces"},
	}

	for _, q := range queues {
		queue, err := ch.QueueDeclare(q.name, true, false, false, false, nil)
		if err != nil {
			ch.Close()
			conn.Close()
			log.Fatal("failed at declare queues:", err)
		}

		err = ch.QueueBind(queue.Name, q.routingKey, "public_exchange", false, nil)
		if err != nil {
			ch.Close()
			conn.Close()
			log.Fatal("failed at bindings queues:", err)
		}
	}

	log.Println("RabbitMQ setup completed successfully")

	return &shared_models.Broker{
		Connection: conn,
		Channel:    ch,
	}
}
