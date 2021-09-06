package main

import (
	"context"
	"fmt"
	"github.com/wenchong-wei/quant-order/adapter"
	"github.com/wenchong-wei/quant-order/pub"
)

func main() {
	adapter.InitClient()
	fmt.Println(adapter.Client.CreateGridOrder(context.TODO(), &pub.CreateGridOrderReq{
		Uid:        1,
		BrokerType: 1,
		QuantType:  1,
		AssetType:  3,
		AssetCode:  "BABA",
		TotalMoney: 300000000,
		GridNum:    10,
		GridMax:    300000,
		GridMin:    0,
	}))

}
