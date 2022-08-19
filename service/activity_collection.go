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

func (s *activityCollection) Create(in model.CreateActivityCollectionReq) (err error) {
	aids := make([]int, 0)
	for _, v := range in.Activities {
		aids = append(aids, v.Id)
	}
	activits := Activity.GetByIds(aids)
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
	aids := make([]int, 0)
	for _, v := range in.Activities {
		aids = append(aids, v.Id)
	}
	activits := Activity.GetByIds(aids)
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
