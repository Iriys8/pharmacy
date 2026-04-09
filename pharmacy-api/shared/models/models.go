package models

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

type RequestContext struct {
	Query   json.RawMessage
	Context json.RawMessage
}

type Message struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Immediate  bool
	Publishing amqp.Publishing
}

func (b *Broker) Close() {
	if b.Channel != nil {
		b.Channel.Close()
	}
	if b.Connection != nil {
		b.Connection.Close()
	}
}
