package main

import (
	"log"

	shared_models "pharmacy-api/shared/models"
)

func setupBroker(broker *shared_models.Broker) (err error) {

	err = broker.Channel.ExchangeDeclare("public_exchange", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to declare exchange:", err)
		return
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
		queue, err := broker.Channel.QueueDeclare(q.name, true, false, false, false, nil)
		if err != nil {
			log.Fatal("failed at declare queues:", err)
			return err
		}

		err = broker.Channel.QueueBind(queue.Name, q.routingKey, "public_exchange", false, nil)
		if err != nil {
			log.Fatal("failed at bindings queues:", err)
			return err
		}
	}

	log.Println("RabbitMQ setup completed successfully")
	return
}
