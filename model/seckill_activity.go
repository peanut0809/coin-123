package model

import (
	"github.com/gogf/gf/os/gtime"
)

type SeckillActivity struct {
	Id                int         `orm:"id,primary" json:"id"` // 活动iD
	Alias             string      `orm:"alias" json:"alias"`   // 别名
	PublisherId       string      `orm:"publisher_id" json:"publisherId"`
	ActivityStartTime *gtime.Time `orm:"activity_start_time" json:"activityStartTime"` // 活动开始时间
	ActivityEndTime   *gtime.Time `orm:"activity_end_time" json:"activityEndTime"`     //活动结束时间
	StartTime         *gtime.Time `orm:"start_time" json:"startTime"`                  // 普通用户认购开始时间
	Price             int         `orm:"price" json:"price"`                           // 发售价,单位：分
	AppId             string      `orm:"app_id" json:"appId"`                          // 应用ID
	AssetType         string      `orm:"asset_type" json:"assetType"`                  // 资产类型
	AssetType2        string      `orm:"asset_type2" json:"assetType2"`                // 资产类型
	AssetType3        string      `orm:"asset_type3" json:"assetType3"`                // 资产类型
	TemplateId        string      `orm:"template_id" json:"templateId"`                // 模板ID
	SumNum            int         `orm:"sum_num" json:"sumNum"`                        // 总数
	RemainNum         int         `orm:"remain_num" json:"remainNum"`
	LimitBuy          int         `orm:"limit_buy" json:"limitBuy"`
	CreatedAt         *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt         *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}

type SeckillActivityFull struct {
	*SeckillActivity
	LastSec int64 `json:"lastSec"`
	Status  int   `json:"status"`
}

const SeckillActivityStatus_Wait_Start = 0
const SeckillActivityStatus_Ing = 1
const SeckillActivityStatus_End = 2
