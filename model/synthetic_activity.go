package model

import "github.com/gogf/gf/os/gtime"

type SyntheticActivity struct {
	Id          int         `orm:"id" json:"id"`
	PublisherId string      `orm:"publisher_id" json:"publisherId"`
	Name        string      `orm:"name" json:"name"`
	AppId       string      `orm:"app_id" json:"appId"`
	AssetType   string      `orm:"asset_type" json:"assetType"`
	TemplateId  string      `orm:"template_id" json:"templateId"`
	Sum         int         `orm:"sum" json:"sum"`
	OutNum      int         `orm:"out_num" json:"outNum"`
	RemainNum   int         `orm:"remain_num" json:"remainNum"`
	Cover       string      `orm:"cover" json:"cover"`
	Rule        string      `orm:"rule" json:"rule"`
	StartTime   *gtime.Time `orm:"start_time" json:"startTime"`
	EndTime     *gtime.Time `orm:"end_time" json:"endTime"`
	Condition   *string     `orm:"condition" json:"condition"`
	CreatedAt   *gtime.Time `orm:"created_at" json:"createdAt"`
	UpdatedAt   *gtime.Time `orm:"updated_at" json:"updatedAt"`
}

type ConditionItem struct {
	AppId      string `json:"appId"`
	AssetType  string `json:"assetType"`
	TemplateId string `json:"templateId"`
	Num        int    `json:"num"`
}

type SyntheticActivityReq struct {
	SyntheticActivity
	ConditionArr []ConditionItem `json:"conditionArr"`
}
