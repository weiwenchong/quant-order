package main

import (
	"context"
	"fmt"
	"github.com/wenchong-wei/quant-order/pub"
)

func main() {
	fmt.Println(pub.Client.CreateGridOrder(context.TODO(), &pub.CreateGridOrderReq{
		Uid:        1,
		BrokerType: 1,
		QuantType:  1,
		AssetType:  1,
		AssetCode:  "513",
		TotalMoney: 1000,
		GridNum:    10,
		GridMax:    1000,
		GridMin:    0,
	}))
}
