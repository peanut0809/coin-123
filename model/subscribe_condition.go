package model

import "github.com/gogf/gf/os/gtime"

type SubscribeCondition struct {
	Id           uint        `orm:"id,primary" json:"id"` // 条件ID
	Aid          int         `orm:"aid" json:"aid"`       // 活动ID
	PublisherId  string      `orm:"publisher_id" json:"publisherId"`
	AppId        string      `orm:"app_id" json:"appId"`                // 应用ID
	AssetType    string      `orm:"asset_type" json:"assetType"`        // 资产类型
	AssetType2   string      `orm:"asset_type2" json:"assetType2"`      // 资产类型
	AssetType3   string      `orm:"asset_type3" json:"assetType3"`      // 资产类型
	TemplateId   string      `orm:"template_id" json:"templateId"`      // 模板ID
	MetaDataRule string      `orm:"meta_data_rule" json:"metaDataRule"` // 属性数组
	BuyNum       int         `orm:"buy_num" json:"buyNum"`              // 可购买数量
	CreatedAt    *gtime.Time `orm:"created_at" json:"createdAt"`        // 新建时间
	UpdatedAt    *gtime.Time `orm:"updated_at" json:"updatedAt"`        // 更新时间
}

type Asset struct {
	AppId      string `json:"appId"`
	TokenId    string `json:"tokenId"`
	AssetPic   string `json:"assetPic"`
	Icon       string `json:"icon"`
	AssetName  string `json:"assetName"`
	TemplateId string `json:"templateId"`
	Type       string `json:"type,omitempty"`
	Data       string `json:"data,omitempty"`
}
