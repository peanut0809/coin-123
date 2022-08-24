package service

import (
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
)

type activityCollectionContent struct {
}

var ActivityCollectionContent = new(activityCollectionContent)

func (s *activityCollectionContent) GetActivityIds(activityCollectionId int) (ret []int, err error) {
	as := make([]model.ActivityCollectionContent, 0)
	err = g.DB().Model("activity_collection_content").Where("activity_collection_id = ?", activityCollectionId).Scan(&as)
	if err != nil {
		return
	}
	for _, v := range as {
		ret = append(ret, v.ActivityId)
	}
	return
}
