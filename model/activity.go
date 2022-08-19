package model

import (
	"github.com/gogf/gf/os/gtime"
)

type Activity struct {
	Id           int         `orm:"id" json:"id"`
	PublisherId  string      `orm:"publisher_id" json:"publisherId"`
	ActivityId   int         `orm:"activity_id" json:"activityId"`
	ActivityType int         `orm:"activity_type" json:"activityType"`
	CreatedAt    *gtime.Time `orm:"created_at" json:"createdAt"`
	UpdatedAt    *gtime.Time `orm:"updated_at" json:"updatedAt"`
}
