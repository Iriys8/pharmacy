package setup

import (
	"fmt"
	"log"
	"os"

	shared_models "pharmacy-api/shared/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectBroker() *shared_models.Broker {
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

	return &shared_models.Broker{
		Connection: conn,
		Channel:    ch,
	}
}
