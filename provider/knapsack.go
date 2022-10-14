package provider

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

const FrozenBySale = "sale"
const FrozenByMhrg = "mhrg"

var KnapsackService = new(knapsack)

type knapsack struct {
}

type getNftInfoRet struct {
	List []KnapsackInfo `json:"list"`
}

type KnapsackInfo struct {
	AppId       string      `json:"appId"`
	TokenId     string      `json:"tokenId"`
	IsFrozen    int         `json:"isFrozen"`
	FrozenBy    string      `json:"frozenBy"`
	Data        string      `json:"data"`
	TemplateId  string      `json:"templateId"`
	Type        string      `orm:"type" json:"type"` // 资产ID
	Type2       string      `orm:"type2" json:"type2"`
	Type3       string      `orm:"type3" json:"type3"`
	AssetName   string      `orm:"asset_name" json:"assetName"`    // 资产名称
	AssetPic    string      `orm:"asset_pic" json:"assetPic"`      // 资产图片URI
	Icon        string      `orm:"icon" json:"icon"`               // 资产背包图标
	Description string      `orm:"description" json:"description"` // 资产描述
	CreatedAt   *gtime.Time `orm:"created_at" json:"createdAt"`
	CanUse      int         `orm:"can_use" json:"canUse"`
	No          int         `orm:"no" json:"no"`
	Total       int         `orm:"total" json:"total"`
	NfrTime     *gtime.Time `json:"nfrTime"`
	PublisherId string      `json:"publisherId"`
	UserId      string      `json:"userId"`
	ChangeNfrYn int         `json:"changeNfrYn"`
}

func (s *knapsack) GetNftInfoByAppIdTokenId(appId, tokenId string) (ret KnapsackInfo, err error) {
	params := &map[string]interface{}{
		"appId":   appId,
		"tokenId": tokenId,
	}
	result, e := utils.SendJsonRpc(context.Background(), "knapsack", "AssetKnapsack.StatusByAppIdAndTokenId", params)
	if e != nil {
		err = e
		return
	}
	err = json.Unmarshal([]byte(gconv.String(result)), &ret)
	if err != nil {
		return
	}
	return
}

func (s *knapsack) GetNftInfo(userId, appId, tokenId string) (ret KnapsackInfo, err error) {
	params := &map[string]interface{}{
		"userId":   userId,
		"appIds":   []string{appId},
		"tokenIds": []string{tokenId},
		"frozenBy": "all",
	}
	result, e := utils.SendJsonRpc(context.Background(), "knapsack", "AssetKnapsack.GetList", params)
	if e != nil {
		err = e
		return
	}
	fmt.Println(gconv.String(result))
	var info getNftInfoRet
	err = json.Unmarshal([]byte(gconv.String(result)), &info)
	if err != nil {
		return
	}
	if len(info.List) == 0 {
		err = fmt.Errorf("资产无效")
		return
	}
	ret = info.List[0]
	return
}

func (s *knapsack) FrozenAsset(userId string, appId string, tokenId string, isFrozen int, frozenBy string) (err error) {
	params := &map[string]interface{}{
		"userId":   userId,
		"appId":    appId,
		"tokenIds": []string{tokenId},
		"isFrozen": isFrozen,
		"frozenBy": frozenBy,
		"opt": map[string]interface{}{
			"optUserId": userId,
			"optType":   "MARKET",
			"optRemark": "市场挂售",
		},
	}
	_, err = utils.SendJsonRpc(context.Background(), "knapsack", "AssetKnapsack.IsFrozenUpdateBatch", params)
	if err != nil {
		g.Log().Error("AssetKnapsack.IsFrozenUpdateBatch err:%v,params:%+v", err, params)
		return
	}
	return
}

type GetListByTemplateRet struct {
	List []struct {
		Id       int    `json:"id"`
		AppId    string `json:"appId"`
		TokenId  string `json:"tokenId"`
		Metadata struct {
			TemplateId string `json:"templateId"`
		} `json:"metadata"`
	} `json:"list"`
}

func (s *knapsack) GetListByTemplate(userId, appId, templateId string) (ret GetListByTemplateRet, err error) {
	params := &map[string]interface{}{
		"userId":      userId,
		"appIds":      []string{appId},
		"templateIds": []string{templateId},
		"pageNum":     1,
		"pageSize":    10000,
	}
	result, e := utils.SendJsonRpc(context.Background(), "knapsack", "AssetKnapsack.GetListByTemplate", params)
	if e != nil {
		err = e
		return
	}
	err = json.Unmarshal([]byte(gconv.String(result)), &ret)
	if err != nil {
		return
	}
	return
}

func (s *knapsack) DeleteByIds(ids []int, optType string, optRemark string) (err error) {
	params := &map[string]interface{}{
		"ids": ids,
		"opt": g.Map{"optType": optType, "optRemark": optRemark},
	}
	_, err = utils.SendJsonRpc(context.Background(), "knapsack", "AssetKnapsack.DeleteByIds", params)
	if err != nil {
		return
	}
	return
}
