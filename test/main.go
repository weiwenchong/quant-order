package main

import (
	"context"
	"fmt"
	"github.com/weiwenchong/quant-order/service/model/cache"
)

func main() {
	fmt.Println(cache.GetAssetPrice(context.TODO(), 1, "600031"))

}
