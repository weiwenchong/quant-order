package quant

import ()

const (
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
