package quant

import (
	"context"
	"encoding/json"
	. "github.com/wenchong-wei/quant-order/pub"
	"github.com/wenchong-wei/quant-order/service/model/dao"
	"log"
)

const (
	// TaskType
	GRID_TASK = 1
)

var (
	ShBegin int64 = 9*3600 + 30*60
	ShEnd   int64 = 15 * 3600
)

func GetTradeTimeByAssetType(asset_type int32) (start, end int64) {
	switch asset_type {
	case 1:
		return ShBegin, ShEnd
	case 2:
		return ShBegin, ShEnd
	default:
		return 0, 0
	}
}

type Task struct {
	Type int32
	Req  string
}

func ConsumerTask(ctx context.Context, msg string) {
	t := &Task{}
	err := json.Unmarshal([]byte(msg), t)
	if err != nil {
		log.Printf("ConsumerTask err:%v", err)
		return
	}
	switch t.Type {
	case GRID_TASK:
		FactoryGriderTask(t.Req).DoTask(ctx)
	default:
		log.Printf("ConsumerTask invalid type:%d", t.Type)
	}
}

type order struct {
}

var OrderMgr *order

func (m *order) CloseOrder(ctx context.Context, oid int64, uid int64) error {
	fun := "order.CloseOrder -->"
	_, err := dao.Update(ctx, dao.DB, dao.ORDER_INFO, map[string]interface{}{"id": oid, "uid": uid}, map[string]interface{}{"status": 0})
	if err != nil {
		log.Printf("%s dao.Update err:%v", fun, err)
		return err
	}
	return nil
}

func (m *order) GetOrdersByUid(ctx context.Context, uid int64) ([]*OrderInfo, error) {
	fun := "order.GetOrdersByUid -->"
	infos := make([]*dao.OrderInfo, 0)
	err := dao.SelectList(ctx, dao.DB, dao.ORDER_INFO, map[string]interface{}{"uid": uid}, &infos)
	if err != nil {
		log.Printf("%s dao.SelectList err:%v", fun, err)
		return nil, err
	}
	rpcInfos := make([]*OrderInfo, 0)
	for _, info := range infos {
		rpcInfos = append(rpcInfos, OrderInfoTrans2Rpc(info))
	}
	return rpcInfos, nil
}

func OrderInfoTrans2Rpc(info *dao.OrderInfo) *OrderInfo {
	rpcInfo := &OrderInfo{
		Id:        info.Id,
		Uid:       info.Uid,
		BrokeType: info.BrokeType,
		QuantType: info.QuantType,
		AssetType: info.AssetType,
		AssetCode: info.AssetCode,
		Total:     info.Total,
		Grids:     nil,
		Hold:      info.Hold,
		Profit:    info.Profit,
		Freeze:    info.Freeze,
		Status:    info.Status,
		Ct:        info.Ct,
	}
	if info.QuantType == 1 {
		json.Unmarshal([]byte(info.Info), &rpcInfo.Grids)
	}
	return rpcInfo
}
