package model

import (
	"github.com/gogf/gf/os/gtime"
)

type WhiteActivity struct {
	Id                uint        `orm:"id" json:"id"`
	PublisherId       string      `orm:"id" json:"publisherId"`
	Name              string      `orm:"name" json:"name"`
	Price             string      `orm:"price" json:"price"`
	LimitType         int         `orm:"limit_type" json:"limitType"`
	LimitBuy          int         `orm:"limit_buy" json:"LimitBuy"`
	Num               int         `orm:"num" json:"num"`
	Status            int         `orm:"status" json:"status"`
	ActivityStartTime *gtime.Time `orm:"activity_start_time" json:"activityStartTime"`
	ActivityEndTime   *gtime.Time `orm:"activity_end_time" json:"activityEndTime"`
	CreatedAt         *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt         *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}

type WhiteActivityList struct {
	List  []*WhiteActivity `json:"list"`
	Total int              `json:"total"`
}
