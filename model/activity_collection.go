package model

import (
	"github.com/gogf/gf/os/gtime"
)

type ActivityCollection struct {
	Id            int         `orm:"id" json:"id"`                         // 主键
	PublisherId   string      `orm:"publisher_id" json:"publisherId"`      // 发行商ID
	Name          string      `orm:"name" json:"name"`                     // 活动名称
	Remark        string      `orm:"remark" json:"remark"`                 // 活动备注
	Intro         string      `orm:"intro" json:"intro"`                   // 活动简介
	Cover         string      `orm:"cover"json:"cover"`                    // 活动图
	Sort          int         `orm:"sort" json:"sort"`                     // 排序
	ShowStartTime *gtime.Time `orm:"show_start_time" json:"showStartTime"` // 展示开始时间
	ShowEndTime   *gtime.Time `orm:"show_end_time" json:"showEndTime"`     // 展示结束时间
	IsShow        int         `orm:"is_show" json:"isShow"`
	CreatedAt     *gtime.Time `orm:"created_at"json:"createdAt"` // 新建时间
	UpdatedAt     *gtime.Time `orm:"updated_at"json:"updatedAt"` // 更新时间
}

type CreateActivityCollectionReq struct {
	ActivityCollection
	Activities []int `json:"activities"`
}

type AdminActivityCollectionDetail struct {
	ActivityCollection
	Activities []Activity `json:"activities"`
}
