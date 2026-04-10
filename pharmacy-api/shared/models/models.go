package models

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
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

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string
	jwt.RegisteredClaims
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

type RefreshClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
