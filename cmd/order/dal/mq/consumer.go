package mq

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/ozline/tiktok/cmd/order/dal/db"
	"github.com/ozline/tiktok/pkg/constants"
)

func StartSeckillConsumer() {
	msgs, err := Channel.Consume(
		constants.SeckillOrderQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		klog.Fatalf("failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var seckillMsg SeckillMessage
			if err := json.Unmarshal(msg.Body, &seckillMsg); err != nil {
				klog.Errorf("failed to unmarshal seckill message: %v", err)
				msg.Nack(false, false)
				continue
			}

			err := processSeckillOrder(context.Background(), &seckillMsg)
			if err != nil {
				klog.Errorf("failed to process seckill order %s: %v", seckillMsg.OrderNo, err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
			klog.Infof("seckill order processed: %s", seckillMsg.OrderNo)
		}
	}()

	klog.Info("seckill order consumer started")
}

func processSeckillOrder(ctx context.Context, msg *SeckillMessage) error {
	return db.UpdateOrderStatus(ctx, msg.OrderNo, constants.OrderStatusPending)
}
