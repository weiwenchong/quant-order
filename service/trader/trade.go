package trader

import (
	"context"
	"errors"
	. "github.com/wenchong-wei/quant-order/pub"
	"log"
)

type trader interface {
	// 返回买入总价
	Buy(ctx context.Context, assetType int32, assetCode string, num int64, price int64) (tradeMoney int64, err error)
	Sell(ctx context.Context, assetType int32, assetCode string, num int64, price int64) (tradelMoney int64, err error)
	Query(ctx context.Context) (assets []*AssetData, err error)
}

func FactoryTrader(t int32, uid int64) trader {
	switch t {
	case 1:
		// 模拟交易
		return &mockTrader{uid: uid}
	default:
		log.Printf("FactoryTrader err: invalid type")
		return &badTrader{uid: uid}
	}
}

type badTrader struct {
	uid int64
}

func (t *badTrader) Buy(ctx context.Context, assetType int32, assetCode string, num int64, price int64) (tradeMoney int64, err error) {
	return 0, errors.New("invalid Trader")
}

func (t *badTrader) Sell(ctx context.Context, assetType int32, assetCode string, num int64, price int64) (tradelMoney int64, err error) {
	return 0, errors.New("invalid Trader")
}

func (t *badTrader) Query(ctx context.Context) (assets []*AssetData, err error) {
	return nil, errors.New("invalid Trader")
}

type mockTrader struct {
	uid int64
}

func (t *mockTrader) Buy(ctx context.Context, assetType int32, assetCode string, num int64, price int64) (tradeMoney int64, err error) {
	return
}

func (t *mockTrader) Sell(ctx context.Context, assetType int32, assetCode string, num int64, price int64) (tradelMoney int64, err error) {
	return
}

func (t *mockTrader) Query(ctx context.Context) (assets []*AssetData, err error) {
	return
}
