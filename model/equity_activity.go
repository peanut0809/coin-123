package model

import (
	"github.com/gogf/gf/os/gtime"
)

const WHITE_ACTIVITY_STATUS1 = 1 // 上架
const WHITE_ACTIVITY_STATUS2 = 2 // 下架

const WHITE_ACTIVITY_LIMIT_TYPE1 = 1 // 每人限购
const WHITE_ACTIVITY_LIMIT_TYPE2 = 2 // 专属限购 白名单用户

type EquityActivity struct {
	Id                int         `orm:"id,primary" json:"id"`                         // 活动iD
	Name              string      `orm:"Name" json:"Name"`                             // 活动名称
	Price             int         `orm:"price" json:"price"`                           // 发售价,单位：分
	ActivityStartTime *gtime.Time `orm:"activity_start_time" json:"activityStartTime"` // 活动开始时间
	ActivityEndTime   *gtime.Time `orm:"activity_end_time" json:"activityEndTime"`     // 活动结束时间
	LimitBuy          int         `orm:"limit_buy" json:"limitBuy"`                    // 限购类型 1 按每人限购 2 白名单限购  1 每人限购数量
	LimitType         int         `orm:"limit_type" json:"limitType"`                  // 限购类型 1 按每人限购 2 白名单限购
	Number            int         `orm:"number" json:"number"`                         // 总数量
	Status            int         `orm:"status" json:"status"`                         // 活动状态1:上架 2:下架
	CreatedAt         *gtime.Time `orm:"created_at" json:"createdAt"`                  // 新建时间
	UpdatedAt         *gtime.Time `orm:"updated_at" json:"updatedAt"`                  // 更新时间
}

type CreateWhiteActivityReq struct {
	EquityActivity
	PriceYuan string `json:"priceYuan"`
}
