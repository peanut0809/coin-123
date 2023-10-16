package model

import "github.com/gogf/gf/os/gtime"

type SeckillWaitPayOrder struct {
	Id          int         `orm:"id" json:"id"`
	OrderNo     string      `orm:"order_no" json:"orderNo"`
	PayExpireAt *gtime.Time `orm:"pay_expire_at" json:"payExpireAt"`
	CreatedAt   *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt   *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}
