package logic

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/wenchong-wei/quant-order/service/model/cache"
	"github.com/wenchong-wei/quant-order/service/quant"
	"log"
)

func consumRedisTopic(ctx context.Context) {
	fun := "consumRedisTopic -->"
	sub := cache.Client.Subscribe("task_1")
	for {
		message, err := sub.Receive()
		if err != nil {
			log.Printf("%s Receive err:%v", fun, err)
			continue
		}
		switch ms := message.(type) {
		case redis.Message:
			quant.FactoryGriderTask(ms.Payload).DoTask(ctx)
			log.Printf("%s receive msg:%s", fun, ms.Payload)
		default:
			log.Printf("%s receive err msg :%v", fun, ms)
		}
		continue
	}
}
