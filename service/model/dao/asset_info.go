package dao

import (
	"context"
	"log"
	"strings"
)

func GetIfNotInsertAssetInfo(ctx context.Context, t int32, code string) (info *AssetInfo, err error) {
	fun := "GetIfNotInsertAssetInfo -->"
	info = new(AssetInfo)
	err = SelectOne(ctx, DB, ASSET_INFO, map[string]interface{}{"type": t, "code": code}, info)
	if err != nil {
		if strings.Contains(err.Error(), "empty result") {
			info = nil
		} else {
			log.Printf("%s SelectOne err:%v", fun, err)
			return nil, err
		}
	}
	if info != nil {
		return info, nil
	} else {
		var buyLimit int32
		if t == 1 || t == 2 || t == 4 {
			// 上 深 港限制100股
			buyLimit = 100
		} else if t == 3 {
			buyLimit = 1
		}
		// todo 填充name
		id, err := Insert(ctx, DB, ASSET_INFO, []map[string]interface{}{{"type": t, "code": code, "buylimit": buyLimit}})
		if err != nil {
			log.Printf("%s Insert err:%v", fun, err)
			return nil, err
		}
		return &AssetInfo{
			Id:       id,
			Type:     t,
			Code:     code,
			Name:     "",
			Buylimit: buyLimit,
		}, nil
	}
}
