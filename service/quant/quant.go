package quant

import (
	"context"
	"encoding/json"
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
