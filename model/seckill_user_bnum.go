package model

import "github.com/gogf/gf/os/gtime"

type SeckillUserBnum struct {
	Id        int         `orm:"id" json:"id"`
	Aid       int         `orm:"aid" json:"aid"`
	UserId    string      `orm:"user_id" json:"userId"`
	CanBuy    int         `orm:"can_buy" json:"canBuy"`
	CreatedAt *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}
