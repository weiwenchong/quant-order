package trader

import (
	"context"
	"errors"
	mocktranderAdapter "github.com/weiwenchong/quant-mocktrander/adapter"
	mocktrander "github.com/weiwenchong/quant-mocktrander/pub"
	. "github.com/weiwenchong/quant-order/pub"
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
	res, err := mocktranderAdapter.Client.Buy(ctx, &mocktrander.BuyReq{
		Uid:   t.uid,
		Type:  assetType,
		Code:  assetCode,
		Num:   num,
		Price: price,
	})
	if err != nil {
		return 0, err
	}
	return res.Total, nil
}

func (t *mockTrader) Sell(ctx context.Context, assetType int32, assetCode string, num int64, price int64) (tradelMoney int64, err error) {
	res, err := mocktranderAdapter.Client.Sell(ctx, &mocktrander.SellReq{
		Uid:   t.uid,
		Type:  assetType,
		Code:  assetCode,
		Num:   num,
		Price: price,
	})
	if err != nil {
		return 0, err
	}
	return res.Total, nil
}

func (t *mockTrader) Query(ctx context.Context) (assets []*Asset, err error) {
	res, err := mocktranderAdapter.Client.Query(ctx, &mocktrander.QueryReq{Uid: t.uid})
	if err != nil {
		return nil, err
	}
	for _, asset := range res.Assets {
		assets = append(assets, &Asset{
			AssetType: asset.Type,
			AssetCode: asset.Code,
			AssetNum:  asset.Num,
		})
	}
	return
}
