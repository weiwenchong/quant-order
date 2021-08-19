package dao

const (
	ASSET_INFO = "order_asset_info"
	ORDER_INFO = "order_order_info"
)

type AssetInfo struct {
	Id       int64  `json:"id" bdb:"id"`
	Type     int32  `json:"type" bdb:"type"`         // 1上证 2深证 3美股 4港股
	Code     string `json:"code" bdb:"code"`         // 股票代码
	Name     string `json:"name" bdb:"name"`         // 股票名
	Buylimit int32  `json:"buylimit" bdb:"buylimit"` // 购买最低数目
}

type OrderInfo struct {
	Id        int64  `json:"id" bdb:"id"`
	Uid       int64  `json:"uid" bdb:"uid"`
	BrokeType int32  `json:"broketype" bdb:"broketype"` // 券商类型 1模拟交易
	QuantType int32  `json:"quanttype" bdb:"quanttype"` // 1网格
	AssetType int32  `json:"assettype" bdb:"assettype"` // 1上证 2深证 3美股 4港股
	AssetCode string `json:"assetcode" bdb:"assetcode"` // 股票代码
	Total     int64  `json:"total" bdb:"total"`         // 投入资金总量
	Info      string `json:"info" bdb:"info"`           // 策略信息
	Hold      int64  `json:"hold" bdb:"hold"`           // 持有数量
	Profit    int64  `json:"profit" bdb:"profit"`       // 盈利
	Freeze    int64  `json:"freeze" bdb:"freeze"`       // 冻结数量，针对t+1
	Ct        int64  `json:"ct" bdb:"ct"`
}
