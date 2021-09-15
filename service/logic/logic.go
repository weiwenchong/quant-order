package logic

import (
	"context"
	. "github.com/weiwenchong/quant-order/pub"
	"github.com/weiwenchong/quant-order/service/quant"
	"log"
)

type GrpcOrder struct {
	//UnimplementedOrderServer
}

func (m *GrpcOrder) CreateGridOrder(ctx context.Context, req *CreateGridOrderReq) (*CreateGridOrderRes, error) {
	fun := "GrpcOrder.CreateOrder -->"
	log.Printf("%s incall", fun)
	grider := quant.FactoryGrider()
	err := grider.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("%s CreateOrder succeed grider:%v", fun, grider)
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

func (m *GrpcOrder) CloseOrder(ctx context.Context, req *CloseOrderReq) (*CloseOrderRes, error) {
	fun := "GrpcOrder.CloseOrder -->"
	log.Printf("%s incall", fun)

	if err := quant.OrderMgr.CloseOrder(ctx, req.Oid, req.Uid); err != nil {
		return nil, err
	}

	log.Printf("%s succeed, req:%v", fun, req)
	return &CloseOrderRes{}, nil
}

func (m *GrpcOrder) GetOrdersByUid(ctx context.Context, req *GetOrdersByUidReq) (*GetOrdersByUidRes, error) {
	fun := "GrpcOrder.GetOrdersByUid -->"
	log.Printf("%s incall", fun)

	orders, err := quant.OrderMgr.GetOrdersByUid(ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	log.Printf("%s succeed, req:%v", fun, err)
	return &GetOrdersByUidRes{Orders: orders}, nil
}

func (m *GrpcOrder) GetGridTrial(ctx context.Context, req *CreateGridOrderReq) (*CreateGridOrderRes, error) {
	fun := "GrpcOrder.GetGridTrial -->"
	log.Printf("%s incall", fun)

	grider := quant.FactoryGrider()
	err := grider.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("%s succeed, req:%v, GridOrder:%v", fun, req, grider)
	return &CreateGridOrderRes{
		TotalMoney: grider.TotalMoney,
		Grids:      grider.Grids,
	}, nil
}
