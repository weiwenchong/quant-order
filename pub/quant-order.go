package pub

import (
	"google.golang.org/grpc"
	"log"
)

const PORT = "172.17.0.3:10001"

var Client OrderClient

func InitClient() {
	log.Printf("initClient start")
	conn, err := grpc.Dial(PORT, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Panicf("InitClient quant-order err:%v", err)
	}
	log.Printf("initClient conn succeed")
	Client = NewOrderClient(conn)
	log.Printf("initClient succeed")
	go func() {
		defer func() {
			log.Printf("conn close start")
			conn.Close()
			log.Printf("conn close")
		}()
		select {}
	}()
	log.Printf("initClient end")
}
