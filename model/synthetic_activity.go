package model

import "github.com/gogf/gf/os/gtime"

type SyntheticActivityDetail struct {
	SyntheticActivity
	ConditionArr    []ConditionItem `json:"conditionArr"`
	AssetCateString string          `json:"assetCateString"`
	AssetTotal      int             `json:"assetTotal"`
	AssetCreateAt   *gtime.Time     `json:"assetCreateAt"`
	AssetDetailImg  string          `json:"assetDetailImg"`
	AssetPic        string          `json:"assetPic"`
	AssetName       string          `json:"assetName"`
	ChainName       string          `json:"chainName"`
	ChainAddr       string          `json:"chainAddr"`
	ChainType       int             `json:"chainType"`
}

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
	Open        int         `orm:"open" json:"open"`
	CreatedAt   *gtime.Time `orm:"created_at" json:"createdAt"`
	UpdatedAt   *gtime.Time `orm:"updated_at" json:"updatedAt"`
}

type ConditionItem struct {
	AppId      string `json:"appId"`
	AssetType  string `json:"assetType"`
	TemplateId string `json:"templateId"`
	Name       string `json:"name"`
	Num        int    `json:"num"`
	Cover      string `json:"cover"`
	TokenId    string `json:"tokenId,omitempty"`
}

type SyntheticActivityReq struct {
	SyntheticActivity
	ConditionArr []ConditionItem `json:"conditionArr"`
}

type SyntheticActivityFull struct {
	SyntheticActivity
	StatusTxt string `json:"statusTxt"`
}

type SyntheticActivityList struct {
	Total int                     `json:"total"`
	List  []SyntheticActivityFull `json:"list"`
}

type SyntheticReq struct {
	UserId       string          `json:"userId"`
	Aid          int             `json:"aid"`
	OrderNo      string          `json:"orderNo"`
	PublisherId  string          `json:"publisherId"`
	ConditionArr []ConditionItem `json:"conditionArr"`
}

type SyntheticRet struct {
	Step    string `json:"step"`
	OrderNo string `json:"orderNo"`
	Reason  string `json:"reason"`
}
