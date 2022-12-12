package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"time"
)

type equity struct{}

var Equity = new(equity)

// 活动列表
func (c *equity) List(publisherId string, pageNum int, pageSize int) (res model.EquityActivityList, err error) {
	var equity []*model.EquityActivity
	m := g.DB().Model("equity_activity")
	if publisherId != "" {
		m = m.Where("publisher_id = ?", publisherId)
	}
	res.Total, err = m.Count()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if res.Total == 0 {
		return
	}
	err = m.Order("id DESC").Page(pageNum, pageSize).Scan(&equity)
	if err != nil {
		return
	}
	timeNow := time.Now()
	for _, v := range equity {
		res.List = append(res.List, v)
		if v.ActivityStartTime.Unix() > timeNow.Unix() {
			v.ActivityStatus = 0
			v.ActivityStatusTxt = "未开始"
		}
		if timeNow.Unix() > v.ActivityStartTime.Unix() && timeNow.Unix() < v.ActivityEndTime.Unix() {
			v.ActivityStatus = 1
			v.ActivityStatusTxt = "进行中"
		}
		if timeNow.Unix() > v.ActivityEndTime.Unix() {
			v.ActivityStatus = 2
			v.ActivityStatusTxt = "已结束"
		}
	}
	return
}

// 活动详情
func (c *equity) Info(activityId int) (res model.EquityActivityFull, err error) {
	m := g.DB().Model("equity_activity")
	m.Where("id = ?", activityId)
	err = m.Scan(&res)
	res.LastSec = res.ActivityStartTime.Unix() - time.Now().Unix()
	if res.LastSec < 0 {
		res.LastSec = 0
	}
	timeNow := time.Now()
	if res.ActivityStartTime.Unix() > timeNow.Unix() {
		res.ActivityStatus = 0
		res.ActivityStatusTxt = "未开始"
	}
	if timeNow.Unix() > res.ActivityStartTime.Unix() && timeNow.Unix() < res.ActivityEndTime.Unix() {
		res.ActivityStatus = 1
		res.ActivityStatusTxt = "进行中"
	}
	if timeNow.Unix() > res.ActivityEndTime.Unix() {
		res.ActivityStatus = 2
		res.ActivityStatusTxt = "已结束"
	}
	if err != nil {
		return
	}
	return
}

// 创建订单
func (c *equity) Create(req model.EquityOrderReq) {
	// 获取活动详情
	activityInfo, e := c.GetValidDetail(req.Id)
	if e != nil {
		EquityOrder.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	// 创建订单
	err := EquityOrder.Create(&req, activityInfo)
	if err != nil {
		EquityOrder.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	return
}

// 获取活动详情
func (c *equity) GetValidDetail(id int) (ret *model.EquityActivity, err error) {
	now := time.Now()
	err = g.DB().Model("equity_activity").Where("id = ?", id).Scan(&ret)
	if err != nil {
		return
	}
	if ret == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	if now.Unix() > ret.ActivityStartTime.Unix() && now.Unix() < ret.ActivityEndTime.Unix() {
		ret.Status = model.EquityActivityStatusIng
	} else {
		if now.Unix() < ret.ActivityStartTime.Unix() {
			ret.Status = model.EquityActivityStatusWait
			err = fmt.Errorf("活动暂未开始")
			//ret.LastSec = as.ActivityStartTime.Unix() - now.Unix()
		} else {
			err = fmt.Errorf("活动已结束")
			ret.Status = model.EquityActivityStatusEnd
		}
	}
	return ret, err
}