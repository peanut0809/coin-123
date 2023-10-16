package provider

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"context"
	"github.com/gogf/gf/util/gconv"
)

type land struct {
}

var Land = new(land)

func (s *land) GetLandCount(landTokenId string, userId string) (count int, err error) {
	params := &map[string]interface{}{
		"landTokenId": landTokenId,
		"userId":      userId,
	}
	result, err := utils.SendJsonRpc(context.Background(), "land", "LandBase.GetLandCount", params)
	if err != nil {
		return
	}
	count = gconv.Int(result)
	return
}
