package model

import (
	"github.com/gogf/gf/os/gtime"
)

type EquityUser struct {
	Id          int         `orm:"id,primary" json:"id"`            // 活动ID
	PublisherId string      `orm:"publisher_id" json:"publisherId"` // 发行商ID
	ActivityId  int         `orm:"activity_id" json:"activityId"`   // 活动ID
	UserId      string      `orm:"userId" json:"userId"`            // 用户ID
	Phone       string      `orm:"phone" json:"phone"`              // 手机号
	LimitNum    int         `orm:"limit_num" json:"limitNum"`       // 限购数量
	Price       int         `orm:"price" json:"price"`              // 发售价,单位：分
	Status      int         `orm:"status" json:"status"`            // 活动状态1:正常 2:失效
	Note        string      `orm:"note" json:"note"`                // 备注失效原因
	CreatedAt   *gtime.Time `orm:"created_at" json:"createdAt"`     // 新建时间
	UpdatedAt   *gtime.Time `orm:"updated_at" json:"updatedAt"`     // 更新时间
}
