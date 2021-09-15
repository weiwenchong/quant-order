package logic

import (
	"context"
	"github.com/weiwenchong/quant-order/service/quant"
)

func InitLogic() {
	ctx := context.TODO()
	quant.UpdateFreeze(ctx)
	go consumRedisTopic(ctx)
}
