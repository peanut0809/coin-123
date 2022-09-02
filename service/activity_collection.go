package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"time"
)

type activityCollection struct {
}

var ActivityCollection = new(activityCollection)

func (s *activityCollection) Delete(id int, publisherId string) (err error) {
	err = g.DB().Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
		r, e := tx.Exec("DELETE FROM activity_collection WHERE id = ? AND publisher_id = ?", id, publisherId)
		if e != nil {
			return e
		}
		affectedNum, _ := r.RowsAffected()
		if affectedNum != 1 {
			return fmt.Errorf("删除失败")
		}
		_, e = tx.Exec("DELETE FROM activity_collection_content WHERE activity_collection_id = ?", id)
		if e != nil {
			return e
		}
		return nil
	})
	return
}

func (s *activityCollection) DetailByClient(publisherId string, pageNum int, pageSize int) (ret model.ClientActivityCollectionList, err error) {
	m := g.DB().Model("activity_collection").Where("publisher_id = ? AND NOW() >= show_start_time AND NOW() <= show_end_time", publisherId)
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	as := make([]model.ActivityCollection, 0)
	err = m.Order("start_time").Page(pageNum, pageSize).Scan(&as)
	if err != nil {
		return
	}
	now := time.Now()
	for _, v := range as {
		item := model.ActivityCollectionFull{
			ActivityCollection: v,
			Status:             0,
			StatusTxt:          "活动未开始",
		}
		if now.Unix() >= item.StartTime.Unix() && now.Unix() <= item.EndTime.Unix() {
			item.Status = 1
			item.StatusTxt = "活动进行中"
		} else {
			if now.Unix() > item.EndTime.Unix() {
				item.Status = 2
				item.StatusTxt = "活动已结束"
			}
		}
		ret.List = append(ret.List, item)
	}
	return
}

func (s *activityCollection) ListByClient(id int, publisherId string, pageNum int, pageSize int) (ret model.ClientActivityCollectionList, err error) {
	m := g.DB().Model("activity_collection").Where("publisher_id = ? AND NOW() >= show_start_time AND NOW() <= show_end_time", publisherId)
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	if id != 0 {
		m = m.Where("id = ?", id)
	}
	as := make([]model.ActivityCollection, 0)
	err = m.Order("start_time").Page(pageNum, pageSize).Scan(&as)
	if err != nil {
		return
	}
	now := time.Now()
	for _, v := range as {
		item := model.ActivityCollectionFull{
			ActivityCollection: v,
			Status:             0,
			StatusTxt:          "活动未开始",
		}
		if now.Unix() >= item.StartTime.Unix() && now.Unix() <= item.EndTime.Unix() {
			item.Status = 1
			item.StatusTxt = "活动进行中"
		} else {
			if now.Unix() > item.EndTime.Unix() {
				item.Status = 2
				item.StatusTxt = "活动已结束"
			}
		}
		ret.List = append(ret.List, item)
	}
	return
}

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
		m = m.Where("(name like ? OR id = ?)", "%"+searchVal+"%", searchVal)
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
		as := Activity.GetByIds(activityIds)

		subIds := make([]int, 0)
		secKillIds := make([]int, 0)
		for _, v := range as {
			if v.ActivityType == 3 {
				secKillIds = append(secKillIds, v.ActivityId)
			} else {
				subIds = append(subIds, v.ActivityId)
			}
		}
		var (
			subAcMap     map[int]model.SubscribeActivity
			secKillAcMap map[int]model.SeckillActivity
		)
		if len(subIds) != 0 {
			subAcMap = SubscribeActivity.GetByIds(subIds)
		}
		if len(secKillIds) != 0 {
			secKillAcMap = SeckillActivity.GetByIds(secKillIds)
		}

		for _, v := range as {
			item := model.AdminActivityCollectionDetailActivity{
				Activity:       v,
				SumNum:         0,
				PriceYuan:      "",
				ActivityType:   "",
				ActivityStatus: "",
			}
			n := time.Now()
			if v.StartTime.Unix() > n.Unix() {
				item.ActivityStatus = "未开始"
			}
			if n.Unix() > v.StartTime.Unix() && n.Unix() < v.EndTime.Unix() {
				item.ActivityStatus = "进行中"
			}
			if n.Unix() > v.EndTime.Unix() {
				item.ActivityStatus = "已结束"
			}
			if v.ActivityType == 1 {
				item.ActivityType = "优先购"
				item.SumNum = subAcMap[v.ActivityId].SumNum
				item.PriceYuan = fmt.Sprintf("%.2f", float64(subAcMap[v.ActivityId].Price)/100)
			}
			if v.ActivityType == 2 {
				item.ActivityType = "普通购"
				item.SumNum = subAcMap[v.ActivityId].SumNum
				item.PriceYuan = fmt.Sprintf("%.2f", float64(subAcMap[v.ActivityId].Price)/100)
			}
			if v.ActivityType == 3 {
				item.ActivityType = "秒杀"
				item.SumNum = secKillAcMap[v.ActivityId].SumNum
				item.PriceYuan = fmt.Sprintf("%.2f", float64(secKillAcMap[v.ActivityId].Price)/100)
			}
			ret.Activities = append(ret.Activities, item)
		}

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
	//var r sql.Result
	_, err = tx.Model("activity_collection").FieldsEx("updated_at").Data(g.Map{
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
	//affectedNum, _ := r.RowsAffected()
	//if affectedNum != 1 {
	//	err = fmt.Errorf("更新失败")
	//	tx.Rollback()
	//	return
	//}
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
