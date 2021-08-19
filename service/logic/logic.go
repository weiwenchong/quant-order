package logic

import (
	"context"
	. "github.com/wenchong-wei/quant-order/pub"
	"github.com/wenchong-wei/quant-order/service/quant"
	"log"
)

type GrpcOrder struct {
	//UnimplementedOrderServer
}

func (m *GrpcOrder) CreateGridOrder(ctx context.Context, req *CreateGridOrderReq) (*CreateGridOrderRes, error) {
	fun := "GrpcOrder.CreateOrder -->"
	grider := quant.FactoryGrider()
	err := grider.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	// 创建任务
	err = grider.InitTask(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("%s succeed, req:%v, GridOrder:%v", fun, req, grider)
	return &CreateGridOrderRes{
		TotalMoney: grider.TotalMoney,
		Grids:      grider.Grids,
	}, nil
}
