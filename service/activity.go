package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"meta_launchpad/model"
)

type activity struct {
}

var Activity = new(activity)

func (s *activity) GetByIds(ids []int) (ret []model.Activity) {
	_ = g.DB().Model("activity").Where("id in (?)", ids).Scan(&ret)
	return
}

func (s *activity) List(pageNum int, startTime, endTime string, activityType int, status, searchVal, publisherId string) (ret model.AdminActivityList, err error) {
	m := g.DB().Model("activity").Where("publisher_id = ?", publisherId)
	if startTime != "" {
		m = m.Where("start_time >= ?", startTime)
	}
	if endTime != "" {
		m = m.Where("end_time <= ?", endTime)
	}
	if activityType != 0 {
		m = m.Where("activity_type = ?", activityType)
	}
	n := gtime.Now()
	if status == "0" { //未开始
		m = m.Where("start_time > ?", n)
	}
	if status == "1" { //进行中
		m = m.Where("? > start_time AND ? < end_time", n, n)
	}
	if status == "2" { //已结束
		m = m.Where("? > end_time", n)
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
	var as []model.Activity
	err = m.Order("id DESC").Page(pageNum, 20).Scan(&as)
	if err != nil {
		return
	}
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
		item := model.AdminActivityFull{
			Activity:           v,
			SumNum:             0,
			Price:              "",
			ActivityTypeString: "",
			ActivityStatus:     "",
			ActivityStatusTxt:  "",
		}
		if v.ActivityType == 3 {
			item.SumNum = secKillAcMap[v.ActivityId].SumNum
			item.Price = fmt.Sprintf("%.2f", float64(secKillAcMap[v.ActivityId].Price)/100)
			item.Cover = secKillAcMap[v.ActivityId].CoverImgUrl
			item.ActivityTypeString = "秒杀"
			item.Alias = secKillAcMap[v.ActivityId].Alias
		} else {
			item.SumNum = subAcMap[v.ActivityId].SumNum
			item.Price = fmt.Sprintf("%.2f", float64(subAcMap[v.ActivityId].Price)/100)
			item.Cover = subAcMap[v.ActivityId].CoverImgUrl
			item.Alias = subAcMap[v.ActivityId].Alias
		}
		if v.ActivityType == 1 {
			item.ActivityTypeString = "优先购"
		}
		if v.ActivityType == 2 {
			item.ActivityTypeString = "普通购"
		}
		if v.StartTime.Unix() > n.Unix() {
			item.ActivityStatus = "0"
			item.ActivityStatusTxt = "未开始"
		}
		if n.Unix() > v.StartTime.Unix() && n.Unix() < v.EndTime.Unix() {
			item.ActivityStatus = "1"
			item.ActivityStatusTxt = "进行中"
		}
		if n.Unix() > v.EndTime.Unix() {
			item.ActivityStatus = "2"
			item.ActivityStatusTxt = "已结束"
		}
		ret.List = append(ret.List, item)
	}
	return
}
