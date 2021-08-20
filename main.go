package main

import (
	. "github.com/wenchong-wei/quant-order/pub"
	"github.com/wenchong-wei/quant-order/service/logic"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Panicf("quant-order service listen err:%v", err)
		return
	}

	logic.InitLogic()

	s := grpc.NewServer()
	RegisterOrderServer(s, new(logic.GrpcOrder))
	log.Println("order start")
	if err = s.Serve(lis); err != nil {
		log.Panicf("quant-order serve err:%v", err)
	}

}
