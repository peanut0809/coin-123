package service

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
)

type activityCollection struct {
}

var ActivityCollection = new(activityCollection)

func (s *activityCollection) List(publisherId string, pageNum int, createStartTime, createEndTime string, showStartTime, showEndTime string, searchVal string, pageSize int) (ret model.AdminActivityCollectionList, err error) {
	m := g.DB().Model("activity_collection").Where("publisher_id = ?", publisherId)
	if createStartTime != "" && createEndTime != "" {
		m = m.Where("created_at >= ? AND created_at <= ?", createStartTime, createEndTime)
	}
	if showStartTime != "" {
		m = m.Where("show_start_time <= ?", showStartTime)
	}
	if showEndTime != "" {
		m = m.Where("show_end_time >= ?", showEndTime)
	}
	if searchVal != "" {
		m = m.Where("name like ?", "%"+searchVal+"%")
	}
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	as := make([]model.ActivityCollection, 0)
	err = m.Order("id DESC").Page(pageNum, pageSize).Scan(&as)
	if err != nil {
		return
	}
	for _, v := range as {
		item := model.ActivityCollectionFull{
			ActivityCollection: v,
		}
		ret.List = append(ret.List, item)
	}
	return
}

func (s *activityCollection) Detail(id int, publisherId string) (ret model.AdminActivityCollectionDetail, err error) {
	var ac *model.ActivityCollection
	err = g.DB().Model("activity_collection").Where("id = ? AND publisher_id = ?", id, publisherId).Scan(&ac)
	if err != nil {
		return
	}
	if ac == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	ret.ActivityCollection = *ac
	var acContent []model.ActivityCollectionContent
	err = g.DB().Model("activity_collection_content").Where("activity_collection_id", id).Scan(&acContent)
	if err != nil {
		return
	}
	activityIds := make([]int, 0)
	for _, v := range acContent {
		activityIds = append(activityIds, v.ActivityId)
	}
	if len(activityIds) != 0 {
		ret.Activities = Activity.GetByIds(activityIds)
	}
	return
}

func (s *activityCollection) Create(in model.CreateActivityCollectionReq) (err error) {
	activits := Activity.GetByIds(in.Activities)
	var tx *gdb.TX
	tx, err = g.DB().Begin()
	if err != nil {
		return
	}
	var r sql.Result
	r, err = tx.Model("activity_collection").Insert(&in.ActivityCollection)
	if err != nil {
		tx.Rollback()
		return
	}
	aCollectionId, _ := r.LastInsertId()
	for _, v := range activits {
		_, err = tx.Model("activity_collection_content").Insert(&model.ActivityCollectionContent{
			ActivityCollectionId: int(aCollectionId),
			ActivityId:           v.Id,
			Aid:                  v.ActivityId,
			ActivityType:         v.ActivityType,
		})
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit()
	return
}

func (s *activityCollection) Update(in model.CreateActivityCollectionReq) (err error) {
	activits := Activity.GetByIds(in.Activities)
	var tx *gdb.TX
	tx, err = g.DB().Begin()
	if err != nil {
		return
	}
	var r sql.Result
	r, err = tx.Model("activity_collection").FieldsEx("updated_at").Data(g.Map{
		"name":            in.Name,
		"remark":          in.Remark,
		"intro":           in.Intro,
		"cover":           in.Cover,
		"sort":            in.Sort,
		"show_start_time": in.ShowStartTime,
		"show_end_time":   in.ShowEndTime,
	}).Where("id = ? AND publisher_id = ?", in.Id, in.PublisherId).Update()
	if err != nil {
		tx.Rollback()
		return
	}
	affectedNum, _ := r.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("更新失败")
		tx.Rollback()
		return
	}
	_, err = tx.Model("activity_collection_content").Where("activity_collection_id", in.Id).Delete()
	if err != nil {
		tx.Rollback()
		return
	}
	for _, v := range activits {
		_, err = tx.Model("activity_collection_content").Insert(&model.ActivityCollectionContent{
			ActivityCollectionId: in.Id,
			ActivityId:           v.Id,
			Aid:                  v.ActivityId,
			ActivityType:         v.ActivityType,
		})
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit()
	return
}
