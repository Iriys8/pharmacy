package main

import (
	"log"

	models "pharmacy-api/shared/models"
)

func setupBroker(broker *models.Broker) (err error) {

	err = broker.Channel.ExchangeDeclare("local_exchange", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to declare exchange:", err)
		return
	}

	queues := []struct {
		name       string
		routingKey string
	}{
		{"local_goods_queue", "goods"},
		{"local_schedule_queue", "schedule"},
		{"local_orders_queue", "orders"},
		{"local_announces_queue", "announces"},
		{"local_users_queue", "users"},
		{"local_roles_queue", "roles"},
	}

	for _, q := range queues {
		queue, err := broker.Channel.QueueDeclare(q.name, true, false, false, false, nil)
		if err != nil {
			log.Fatal("failed at declare queues:", err)
			return err
		}

		err = broker.Channel.QueueBind(queue.Name, q.routingKey, "local_exchange", false, nil)
		if err != nil {
			log.Fatal("failed at bindings queues:", err)
			return err
		}
	}

	log.Println("RabbitMQ setup completed successfully")
	return
}
