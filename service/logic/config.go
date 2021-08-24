package logic

import (
	"context"
	"github.com/wenchong-wei/quant-order/service/quant"
	task "github.com/wenchong-wei/quant-task/pub"
)

func InitLogic() {
	// init调用rpc
	task.InitClient()

	ctx := context.TODO()
	quant.UpdateFreeze(ctx)
	go consumRedisTopic(ctx)
}
