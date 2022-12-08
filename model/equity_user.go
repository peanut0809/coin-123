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
	Status      int         `orm:"status" json:"status"`            // 活动状态1:正常 2:失效
	Note        string      `orm:"note" json:"note"`                // 备注失效原因
	CreatedAt   *gtime.Time `orm:"created_at" json:"createdAt"`     // 新建时间
	UpdatedAt   *gtime.Time `orm:"updated_at" json:"updatedAt"`     // 更新时间
}

type EquityUserReq struct {
	PublisherId string `json:"publisherId"`
	EquityId    int    `json:"equityId"`
	Page        int    `json:"pageNum"`
	PageSize    int    `json:"pageSize"`
	Phone       int    `json:"phone"`
	Status      int    `json:"status"`
}

type EquityUserFull struct {
	Total int          `json:"total"`
	List  []EquityUser `json:"list"`
}
