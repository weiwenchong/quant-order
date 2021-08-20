package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"log"

	//"github.com/gomodule/redigo/redis"
	"strconv"
)

//var pool *redis.Pool

var Client *redis.Client

func init() {
	log.Printf("cache init")
	//pool = &redis.Pool{
	//	Dial: func() (redis.Conn, error) {
	//		conn, err := redis.Dial("tcp", "172.17.0.1:6379")
	//		if err != nil {
	//			return nil, err
	//		}
	//		_, err = conn.Do("AUTH", "Wwcwwc123")
	//		if err != nil {
	//			return nil, err
	//		}
	//		return conn, nil
	//	},
	//	TestOnBorrow:    func(c redis.Conn, t time.Time) error {
	//		_, err := c.Do("PING")
	//		return err
	//	},
	//	MaxIdle:         2,
	//	MaxActive:       10,
	//	IdleTimeout:     ,
	//	Wait:            false,
	//	MaxConnLifetime: 0,
	//}
	Client = redis.NewClient(&redis.Options{
		Addr:     "172.17.0.1:6379",
		Password: "Wwcwwc123",
		DB:       0,
	})
	log.Printf("cache init succeed")
}

func GetAssetPrice(ctx context.Context, t int32, code string) (price int64, err error) {
	res := Client.Get(fmt.Sprintf("%d_%s", t, code))
	if res.Err() != nil {
		return 0, err
	}
	price, err = strconv.ParseInt(res.Val(), 10, 64)
	return
}
