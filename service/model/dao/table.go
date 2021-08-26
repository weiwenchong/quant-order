package dao

const (
	ASSET_INFO = "order_asset_info"
	ORDER_INFO = "order_order_info"
)

type AssetInfo struct {
	Id       int64  `json:"id" ddb:"id"`
	Type     int32  `json:"type" ddb:"type"`         // 1上证 2深证 3美股 4港股
	Code     string `json:"code" ddb:"code"`         // 股票代码
	Name     string `json:"name" ddb:"name"`         // 股票名
	Buylimit int32  `json:"buylimit" ddb:"buylimit"` // 购买最低数目
}

type OrderInfo struct {
	Id        int64  `json:"id" ddb:"id"`
	Uid       int64  `json:"uid" ddb:"uid"`
	Status    int32  `json:"status" ddb:"status"`       // 订单状态 1正在运行 0已关闭
	BrokeType int32  `json:"broketype" ddb:"broketype"` // 券商类型 1模拟交易
	QuantType int32  `json:"quanttype" ddb:"quanttype"` // 1网格
	AssetType int32  `json:"assettype" ddb:"assettype"` // 1上证 2深证 3美股 4港股
	AssetCode string `json:"assetcode" ddb:"assetcode"` // 股票代码
	Total     int64  `json:"total" ddb:"total"`         // 投入资金总量
	Info      string `json:"info" ddb:"info"`           // 策略信息
	Hold      int64  `json:"hold" ddb:"hold"`           // 持有数量
	Profit    int64  `json:"profit" ddb:"profit"`       // 盈利
	Freeze    int64  `json:"freeze" ddb:"freeze"`       // 冻结数量，针对t+1
	Ct        int64  `json:"ct" ddb:"ct"`
}
