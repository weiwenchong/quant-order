package main

import (
	. "github.com/weiwenchong/quant-order/pub"
	"github.com/weiwenchong/quant-order/service/logic"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	log.Println("service start")
	lis, err := net.Listen("tcp", "0.0.0.0:10001")
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
