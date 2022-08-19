package model

type ActivityCollectionContent struct {
	Id                   int `json:"id"`
	ActivityCollectionId int `orm:"activity_collection_id" json:"activityCollectionId"`
	ActivityId           int `orm:"activity_id" json:"activityId"`
	Aid                  int `orm:"aid" json:"aid"`
	ActivityType         int `orm:"activity_type" json:"activityType"`
}
