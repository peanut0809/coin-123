package model

import (
	"github.com/gogf/gf/os/gtime"
)

const ACTIVITY_TYPE_1 = 1
const ACTIVITY_TYPE_2 = 2
const ACTIVITY_TYPE_3 = 3
const ACTIVITY_TYPE_4 = 4

type Activity struct {
	Id           int         `orm:"id" json:"id"`
	PublisherId  string      `orm:"publisher_id" json:"publisherId"`
	ActivityId   int         `orm:"activity_id" json:"activityId"`
	ActivityType int         `orm:"activity_type" json:"activityType"`
	Name         string      `orm:"name" json:"name"`
	StartTime    *gtime.Time `orm:"start_time" json:"startTime"`
	EndTime      *gtime.Time `orm:"end_time" json:"endTime"`
	CreatedAt    *gtime.Time `orm:"created_at" json:"createdAt"`
	UpdatedAt    *gtime.Time `orm:"updated_at" json:"updatedAt"`
}

type AdminActivityFull struct {
	Activity
	SumNum             int    `json:"sumNum"`
	Price              string `json:"price"`
	ActivityTypeString string `json:"activityTypeString"`
	ActivityStatus     string `json:"activityStatus"`
	ActivityStatusTxt  string `json:"activityStatusTxt"`
	Cover              string `json:"cover"`
	Alias              string `json:"alias"`
	PublisherName      string `json:"publisherName"`
	PublisherUserId    string `json:"publisherUserId"`
	PublisherIcon      string `json:"publisherIcon"`
}

type AdminActivityList struct {
	List  []AdminActivityFull `json:"list"`
	Total int                 `json:"total"`
}
