package logic

import (
	"context"
	"github.com/weiwenchong/quant-order/service/model/cache"
	"github.com/weiwenchong/quant-order/service/quant"
	"log"
)

func consumRedisTopic(ctx context.Context) {
	fun := "consumRedisTopic -->"
	sub := cache.Client.Subscribe("task_1")
	for {
		select {
		case msg := <-sub.Channel():
			log.Printf("%s recive %v", fun, msg)
			quant.ConsumerTask(ctx, msg.Payload)
		}
	}
}
