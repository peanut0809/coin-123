package model

import (
	"github.com/gogf/gf/os/gtime"
)

type AirDropActivityRecord struct {
	Id           int         `orm:"id" json:"id"`                      // 主键
	UserId       string      `orm:"user_id" json:"userId"`             // 用户id
	Phone        string      `orm:"phone" json:"phone"`                // 用户手机号
	ActivityId   int         `orm:"activity_id" json:"activityId"`     // 主表活动id
	Number       int         `orm:"number" json:"number"`              // 空投数量
	ActivityType string      `orm:"activity_type" json:"activityType"` // 空投活动类型 crystal-元晶 speed 加速次数
	Status       int         `orm:"status" json:"status"`              // 状态 1 正常 2 异常
	Message      string      `orm:"message" json:"message"`            // 异常信息记录
	CreatedAt    *gtime.Time `orm:"created_at" json:"createdAt"`       // 新建时间
	UpdatedAt    *gtime.Time `orm:"updated_at" json:"updatedAt"`       // 更新时间
}
