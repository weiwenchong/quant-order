package main

import (
	"context"
	"fmt"
	"github.com/wenchong-wei/quant-order/pub"
)

func main() {
	pub.InitClient()
	fmt.Println(pub.Client.CreateGridOrder(context.TODO(), &pub.CreateGridOrderReq{
		Uid:        1,
		BrokerType: 1,
		QuantType:  1,
		AssetType:  1,
		AssetCode:  "513",
		TotalMoney: 10000000,
		GridNum:    10,
		GridMax:    1000,
		GridMin:    0,
	}))
}
