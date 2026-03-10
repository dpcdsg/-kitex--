package mq

import (
	"fmt"

	"github.com/ozline/tiktok/config"
	"github.com/ozline/tiktok/pkg/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

func Init() {
	url := fmt.Sprintf("amqp://%s:%s@%s/", config.RabbitMQ.Username, config.RabbitMQ.Password, config.RabbitMQ.Addr)

	var err error
	Conn, err = amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		panic(err)
	}

	_, err = Channel.QueueDeclare(
		constants.SeckillOrderQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
}
