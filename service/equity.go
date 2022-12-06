package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/model"
	"time"
)

type equity struct{}

var Equity = new(equity)

// List 活动列表
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
	for _, v := range equity {
		res.List = append(res.List, v)
	}
	return
}

// Info 活动详情
func (c *equity) Info(activityId int) (res model.EquityActivity, err error) {
	m := g.DB().Model("equity_activity")
	m.Where("id = ?", activityId)
	err = m.Scan(&res)
	if err != nil {
		return
	}
	return
}

// Create 下单
func (c *equity) Create(in model.EquityOrderReq) {
	activityInfo, e := c.GetValidDetail(in.Id, in.PublisherId)
	if e != nil {
		c.SetSubResult(model.DoSubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}
	fmt.Println(activityInfo)
	return
}

func (c *equity) SetSubResult(in model.DoSubResult) {
	_, err := g.Redis().Do("SET", fmt.Sprintf(model.SubSetEquityResultKey, in.OrderNo), gconv.String(in), "ex", 3600)
	if err != nil {
		g.Log().Errorf("EquityBuy err:%v", err)
		return
	}
	return
}

func (c *equity) GetValidDetail(id int, publisherId string) (ret model.EquityActivity, err error) {
	var as *model.EquityActivity
	now := time.Now()
	err = g.DB().Model("equity_activity").Where("id = ?", id).Scan(&as)
	if err != nil {
		return
	}
	if as == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	if now.Unix() > as.ActivityStartTime.Unix() && now.Unix() < as.ActivityEndTime.Unix() {
		ret.Status = model.EquityActivityStatusIng
	} else {
		if now.Unix() < as.ActivityStartTime.Unix() {
			ret.Status = model.EquityActivityStatusWait
			//ret.LastSec = as.ActivityStartTime.Unix() - now.Unix()
		} else {
			ret.Status = model.EquityActivityStatusEnd
		}
	}
	return ret, err
}
