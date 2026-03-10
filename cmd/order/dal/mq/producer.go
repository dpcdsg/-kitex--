package mq

import (
	"context"
	"encoding/json"

	"github.com/ozline/tiktok/pkg/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SeckillMessage struct {
	UserId     int64  `json:"user_id"`
	ActivityId int64  `json:"activity_id"`
	OrderNo    string `json:"order_no"`
}

func PublishSeckillOrder(ctx context.Context, msg *SeckillMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return Channel.PublishWithContext(ctx,
		"",
		constants.SeckillOrderQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
