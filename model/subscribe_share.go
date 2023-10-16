package model

import "github.com/gogf/gf/os/gtime"

type SubscribeShare struct {
	Id          int         `orm:"id,primary" json:"id"`            // 条件ID
	SubscribeId int         `orm:"subscribe_id" json:"subscribeId"` // 活动Id
	UserId      string      `orm:"user_id" json:"userId"`           // 用户Id
	CreatedAt   *gtime.Time `orm:"created_at" json:"createdAt"`     // 创建时间
}

type SubscribeShareUpload struct {
	// UserId      string `json:"userId"`
	Alias       string `json:"alias"`
	PublisherId string `json:"publisherId"`
}
