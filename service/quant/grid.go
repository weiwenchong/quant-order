package quant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/os/gcron"
	. "github.com/weiwenchong/quant-order/pub"
	"github.com/weiwenchong/quant-order/service/model/cache"
	"github.com/weiwenchong/quant-order/service/model/dao"
	"github.com/weiwenchong/quant-order/service/trader"
	"github.com/weiwenchong/quant-order/service/util"
	taskAdapter "github.com/weiwenchong/quant-task/adapter"
	task "github.com/weiwenchong/quant-task/pub"
	"log"
	"time"
)

type grider struct {
	order
	Uid        int64
	BrokeType  int32
	AssetType  int32
	AssetCode  string
	TotalMoney int64
	Grids      []*GridData
	Grid       *GridData
}

type griderTask struct {
	Uid       int64
	Id        int64
	BuyPrice  int64
	BrokeType int32
	AssetType int32
	AssetCode string
	Grid      *GridData
	// 0买单 1卖单
	TradeType int32
}

func FactoryGrider() *grider {
	return &grider{}
}

func FactoryGriderTask(message string) *griderTask {
	g := &griderTask{}
	err := json.Unmarshal([]byte(message), g)
	if err != nil {
		log.Printf("FactoryGriderTask json.Unmarshal err:%v", err)
	}
	return g
}

func (m *grider) CreateOrder(ctx context.Context, req *CreateGridOrderReq) (err error) {
	//fun := "grider.CreateOrder -->"
	m.Uid = req.Uid
	m.AssetType = req.AssetType
	m.AssetCode = req.AssetCode
	m.BrokeType = req.BrokerType
	info, err := dao.GetIfNotInsertAssetInfo(ctx, req.AssetType, req.AssetCode)
	if err != nil {
		return
	}
	log.Println(info)
	buyLimit := info.Buylimit
	if buyLimit == 0 {
		buyLimit = 1
	}

	gridDiff := (req.GridMax - req.GridMin) / req.GridNum
	var perMoney int64 = 0
	for i := 0; i < int(req.GridNum); i++ {
		m.Grids = append(m.Grids, &GridData{
			AssetNum: 0,
			GridMax:  req.GridMin + (int64(i+1) * gridDiff),
			GridMin:  req.GridMin + (int64(i) * gridDiff),
		})
		perMoney += req.GridMin + (int64(i) * gridDiff)
	}
	perNumBeforeInt := req.TotalMoney / perMoney
	perNum := perNumBeforeInt - perNumBeforeInt%int64(buyLimit)
	if perNum == 0 {
		return errors.New(fmt.Sprintf("totalmoney too less, perNum:%v", perNumBeforeInt))
	}
	for _, g := range m.Grids {
		g.AssetNum = perNum
	}
	m.TotalMoney = perMoney * perNum

	return nil
}

func (m *grider) InitTask(ctx context.Context) error {
	fun := "grider.InitTask -->"
	info, err := json.Marshal(m.Grids)
	if err != nil {
		return err
	}

	// 购买
	price, err := cache.GetAssetPrice(ctx, m.AssetType, m.AssetCode)
	if err != nil {
		return err
	}
	var (
		buyNum int64 = 0
	)
	minTask := make([]*GridData, 0)
	for _, g := range m.Grids {
		if g.GridMin > price {
			buyNum += g.AssetNum
		} else {
			minTask = append(minTask, g)
		}
	}

	// 网格交易初始化，买入初始资金
	tradMoney, err := trader.FactoryTrader(m.BrokeType, m.Uid).Buy(ctx, m.AssetType, m.AssetCode, buyNum, -1)
	if err != nil {
		return err
	}
	for _, t := range minTask {
		_, err = trader.FactoryTrader(m.BrokeType, m.Uid).Buy(ctx, m.AssetType, m.AssetCode, t.AssetNum, t.GridMin)
		if err != nil {
			log.Printf("%s FactoryTrader minTask err:%v")
		}
	}
	perPrice := tradMoney / int64(len(m.Grids)-len(minTask))
	//var freeze int64
	//if m.AssetType == 1 || m.AssetType == 2 {
	//	freeze = tradMoney
	//}
	id, err := dao.Insert(ctx, dao.DB, dao.ORDER_INFO, []map[string]interface{}{{
		"uid":       m.Uid,
		"status":    1,
		"broketype": m.BrokeType,
		"quanttype": 1,
		"assettype": m.AssetType,
		"assetcode": m.AssetCode,
		"total":     m.TotalMoney,
		"info":      string(info),
		"hold":      0,
		"profit":    0,
		"freeze":    0,
		"ct":        time.Now().Unix(),
	}})

	if err != nil {
		log.Printf("%s dao.Insert err:%v", fun, err)
		return err
	}

	// 批量发送task
	createTaskReq := &task.CreatePriceTaskReq{
		Source: task.SourceService_ORDER,
		Tasks:  make([]*task.PriceTask, 0),
	}
	// 挂价格小于最小价格的卖
	for _, t := range m.Grids {
		tm := &griderTask{
			Uid:       m.Uid,
			Id:        id,
			BuyPrice:  perPrice,
			BrokeType: m.BrokeType,
			AssetType: m.AssetType,
			AssetCode: m.AssetCode,
			Grid:      t,
			// 卖单
			TradeType: 1,
		}
		b, _ := json.Marshal(tm)
		tq, _ := json.Marshal(Task{
			Type: GRID_TASK,
			Req:  string(b),
		})
		createTaskReq.Tasks = append(createTaskReq.Tasks, &task.PriceTask{
			AssetType: m.AssetType,
			AssetCode: m.AssetCode,
			Condition: task.PriceCondition_LESS,
			Price:     t.GridMin,
			Message:   string(tq),
			Uid:       m.Uid,
		})
	}
	_, err = taskAdapter.Client.CreatePriceTask(ctx, createTaskReq)
	if err != nil {
		log.Printf("%s CreatePriceTask err:%v", fun, err)
		return err
	}

	return nil
}

func (m *griderTask) DoTask(ctx context.Context) {
	fun := "griderTask.DoTask -->"
	switch m.TradeType {
	case 0:
		info := new(dao.OrderInfo)
		err := dao.SelectOne(ctx, dao.DB, dao.ORDER_INFO, map[string]interface{}{"id": m.Id}, info)
		if err != nil {
			log.Printf("%s SelectOne err:%v", fun, err)
		}
		if info.Status == 0 {
			return
		}
		// 已经sell，下buy单
		// 更新收益
		_, err = dao.DB.ExecContext(ctx, fmt.Sprintf("update %s set profit = profit + %d, hold = hold - %d where id=%d", dao.ORDER_INFO, (m.Grid.GridMax-m.BuyPrice)*m.Grid.AssetNum, m.Grid.AssetNum, m.Id))
		if err != nil {
			log.Printf("%s dao.DB.ExecContext err:%v", fun, err)
		}
		_, err = trader.FactoryTrader(m.BrokeType, m.Uid).Buy(ctx, m.AssetType, m.AssetCode, m.Grid.AssetNum, m.Grid.GridMin)
		if err != nil {
			log.Printf("%s Buy err:%v", fun, err)
			//return
		}

		m.TradeType = 1
		m.BuyPrice = m.Grid.GridMin
		tm, _ := json.Marshal(m)
		tq, _ := json.Marshal(Task{
			Type: GRID_TASK,
			Req:  string(tm),
		})
		// 买到以后的卖单下任务
		createTaskReq := &task.CreatePriceTaskReq{
			Source: task.SourceService_ORDER,
			Tasks: []*task.PriceTask{{
				AssetType: m.AssetType,
				AssetCode: m.AssetCode,
				Condition: task.PriceCondition_LESS,
				Price:     m.Grid.GridMin,
				Message:   string(tq),
			}},
		}
		_, err = taskAdapter.Client.CreatePriceTask(ctx, createTaskReq)
		if err != nil {
			log.Printf("%s CreatePriceTask err:%v", fun, err)
		}
		return
	case 1:
		// sell
		info := new(dao.OrderInfo)
		err := dao.SelectOne(ctx, dao.DB, dao.ORDER_INFO, map[string]interface{}{}, info)
		if err != nil {
			log.Printf("%s SelectOne err:%v", fun, err)
		}
		if info.Status == 0 {
			return
		}
		var startTime int64
		var freeze int64
		if (info.AssetType == 1 || info.AssetType == 2) && info.Hold-info.Freeze < m.Grid.AssetNum {
			// t+1 今天冻结，明天开始卖
			startTime, _ = GetTradeTimeByAssetType(m.AssetType)
			startTime += util.DayBeginStamp(time.Now().Unix())
		}
		if info.AssetType == 1 || info.AssetType == 2 {
			freeze = m.Grid.AssetNum
		}

		// 已买入，更新持有
		_, err = dao.DB.ExecContext(ctx, fmt.Sprintf("update %s set hold = hold + %d, freeze = freeze + %d where id=%d", dao.ORDER_INFO, m.Grid.AssetNum, freeze, m.Id))
		if err != nil {
			log.Printf("%s dao.DB.ExecContext err:%v", fun, err)
		}

		_, err = trader.FactoryTrader(m.BrokeType, m.Uid).Sell(ctx, m.AssetType, m.AssetCode, m.Grid.AssetNum, m.Grid.GridMax)
		if err != nil {
			log.Printf("%s Sell err:%v", fun, err)
		}
		m.TradeType = 0
		tm, _ := json.Marshal(m)
		tq, _ := json.Marshal(Task{
			Type: GRID_TASK,
			Req:  string(tm),
		})
		createTaskReq := &task.CreatePriceTaskReq{
			Source: task.SourceService_ORDER,
			Tasks: []*task.PriceTask{{
				AssetType: m.AssetType,
				AssetCode: m.AssetCode,
				Condition: task.PriceCondition_GREATER,
				Price:     m.Grid.GridMax,
				StartTime: startTime,
				Message:   string(tq),
			}},
		}
		_, err = taskAdapter.Client.CreatePriceTask(ctx, createTaskReq)
		if err != nil {
			log.Printf("%s CreatePriceTask err:%v", fun, err)
		}
		return
	default:
		log.Printf("%s invalid TaskType:%d", fun, m.TradeType)
		return
	}
}

func UpdateFreeze(ctx context.Context) {
	fun := "UpdateFreeze -->"
	gcron.Add("0 0 0 * * ?", func() {
		_, err := dao.Update(ctx, dao.DB, dao.ORDER_INFO, map[string]interface{}{}, map[string]interface{}{"freeze": 0})
		if err != nil {
			log.Printf("%s err:%v", fun, err)
		} else {
			log.Printf("%s succeed", fun)
		}
	})
}
