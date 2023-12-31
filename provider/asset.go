package provider

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"context"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

type asset struct {
}

var Asset = new(asset)

type assetDestroyInfoRet struct {
	Id         int    `orm:"id" json:"id"`
	AppId      string `orm:"app_id" json:"appId"`
	TemplateId string `orm:"template_id" json:"templateId"`
	Num        int    `orm:"num" json:"num"`
}

//获取销毁量信息
func (s *asset) GetAssetDestroyInfo(appId string, templateId string) (ret assetDestroyInfoRet, err error) {
	params := &map[string]interface{}{
		"appId":      appId,
		"templateId": templateId,
	}
	result, err := utils.SendJsonRpc(context.Background(), "asset", "Asset.GetAssetDestroyInfo", params)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(gconv.String(result)), &ret)
	if err != nil {
		return
	}
	return
}

//获取资产详情
type assetDetailRet struct {
	Total          int         `json:"total"`
	No             int         `json:"no"`
	TokenId        string      `json:"tokenId"`
	TemplateId     string      `json:"templateId"`
	AssetName      string      `json:"assetName"`
	Icon           string      `json:"icon"`
	AssetPic       string      `json:"assetPic"`
	DisposeNum     int         `json:"disposeNum"`
	MetaDataFormat interface{} `json:"metaDataFormat"`
	Description    string      `json:"description"` // 资产描述
	CanUse         int         `json:"canUse"`
}

func (s *asset) GetAssetDetail(appId string, tokenId string) (ret assetDetailRet, err error) {
	params := &map[string]interface{}{
		"appId":   appId,
		"tokenId": tokenId,
	}
	result, err := utils.SendJsonRpc(context.Background(), "asset", "Asset.GetDetail", params)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(gconv.String(result)), &ret)
	if err != nil {
		return
	}
	return
}

type GetMateDataByAksRet struct {
	No         int         `json:"no"`
	Total      int         `json:"total"`
	AssetName  string      `json:"assetName"`
	AssetPic   string      `json:"assetPic"`
	Icon       string      `json:"icon"`
	AppId      string      `json:"appId"`
	TokenId    string      `json:"tokenId"`
	CreateTime *gtime.Time `json:"createTime"`
}

func (s *asset) GetMateDataByAks(appIds []string, tokenIds []string) (ret []GetMateDataByAksRet, retMap map[string]GetMateDataByAksRet, err error) {
	params := &map[string]interface{}{
		"appIds":   appIds,
		"tokenIds": tokenIds,
	}
	err = utils.SendJsonRpcScan(context.Background(), "asset", "Asset.GetMateDataByAks", params, &ret)
	if err != nil {
		g.Log().Errorf("rpcGetMateDataByAks err:%v", err)
		return
	}
	retMap = make(map[string]GetMateDataByAksRet)
	for _, v := range ret {
		retMap[v.AppId+v.TokenId] = v
	}
	return
}

func (s *asset) PublishAssetWithTemplateId(params *map[string]interface{}) (err error) {
	var ret interface{}
	err = utils.SendJsonRpcScan(context.Background(), "asset", "Asset.PublishAssetWithTemplateId", params, &ret)
	if err != nil {
		g.Log().Errorf("PublishAssetWithTemplateId err:%v", err)
		return
	}
	return
}

func (s *asset) GetMateDataByAm(params *map[string]interface{}) (ret GetMateDataByAksRet, err error) {
	err = utils.SendJsonRpcScan(context.Background(), "asset", "Asset.GetMateDataByTemp", params, &ret)
	if err != nil {
		g.Log().Errorf("GetMateDataByTemp err:%v", err)
		return
	}
	return
}

type GetMateDataByTplsItem struct {
	AssetName string `orm:"asset_name" json:"assetName"` // 资产名称
	Icon      string `orm:"icon" json:"icon"`            // 资产背包图标
	AssetPic  string `orm:"asset_pic" json:"assetPic"`
}

func (s *asset) GetMateDataByTpls(params *map[string]interface{}) (ret map[string]GetMateDataByTplsItem, err error) {
	err = utils.SendJsonRpcScan(context.Background(), "asset", "Asset.GetMateDataByTpls", params, &ret)
	if err != nil {
		g.Log().Errorf("GetMateDataByTpls err:%v", err)
		return
	}
	return
}

type GetCanUsedAssetsByTemplateRet struct {
	TokenId string `json:"tokenId"` // tokenId
	AppId   string `orm:"app_id" json:"appId"`
}

func (s *asset) GetCanUsedAssetsByTemplate(params *map[string]interface{}) (ret []GetCanUsedAssetsByTemplateRet, err error) {
	err = utils.SendJsonRpcScan(context.Background(), "asset", "Asset.GetCanUsedAssetsByTemplate", params, &ret)
	if err != nil {
		g.Log().Errorf("GetCanUsedAssetsByTemplate err:%v", err)
		return
	}
	return
}
