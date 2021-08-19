package pub

import (
	"google.golang.org/grpc"
	"log"
)

const PORT = ":10001"

var Client OrderClient

func init() {
	conn, err := grpc.Dial(PORT, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Panicf("InitClient quant-order err:%v", err)
	}
	go func() {
		defer func() {
			log.Printf("conn close start")
			conn.Close()
			log.Printf("conn close")
		}()
		select{}
	}()
	Client = NewOrderClient(conn)
}
