package service

import (
	"meta_launchpad/model"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

type adminEquity struct {
}

var AdminEquity = new(adminEquity)

func (s *adminEquity) Create(in model.EquityActivity) (err error) {

	var tx *gdb.TX
	tx, err = g.DB().Begin()
	if err != nil {
		return
	}
	item, err := tx.Model("equity_activity").Insert(&in)
	if err != nil {
		tx.Rollback()
		return
	}
	activityId, err := item.LastInsertId()
	if err != nil {
		tx.Rollback()
		return
	}
	_, err = tx.Model("activity").Insert(g.Map{
		"name":          in.Name,
		"start_time":    in.ActivityStartTime,
		"end_time":      in.ActivityEndTime,
		"publisher_id":  in.PublisherId,
		"activity_id":   activityId,
		"activity_type": model.ACTIVITY_TYPE_4,
	})
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}
