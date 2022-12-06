package service

import "C"
import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
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
	for _, v := range equity {
		res.List = append(res.List, v)
	}
	return
}

// 活动详情
func (c *equity) Info(activityId int) (res model.EquityActivity, err error) {
	m := g.DB().Model("equity_activity")
	m.Where("id = ?", activityId)
	err = m.Scan(&res)
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
		c.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	if activityInfo.Status == model.EquityActivityStatusEnd {
		c.SetSubResult(model.EquitySubResult{
			Reason:  "活动已结束",
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	var tx *gdb.TX
	tx, e = g.DB().Begin()
	if e != nil {
		c.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	// 扣除库存
	r, e := tx.Exec("UPDATE equity_activity SET number = number - ? WHERE id = ?", req.Num, activityInfo.Id)
	if e != nil {
		tx.Rollback()
		c.SetSubResult(model.EquitySubResult{
			Reason:  "库存不足",
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	affectedNum, _ := r.RowsAffected()
	if affectedNum != 1 {
		tx.Rollback()
		c.SetSubResult(model.EquitySubResult{
			Reason:  "更新库存失败",
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	limitNum := 0
	// 判断白名单
	if activityInfo.LimitType == model.EQUITY_ACTIVITY_LIMIT_TYPE2 {
		var user model.EquityUser
		err := g.DB().Model("equity_user").
			Where("activity_id = ?", req.Id).
			Where("user_id = ?", req.UserId).
			Scan(&user)
		if err != nil {
			tx.Rollback()
			c.SetSubResult(model.EquitySubResult{
				Reason:  e.Error(),
				Step:    "fail",
				OrderNo: req.OrderNo,
			})
			return
		}
		if user.UserId == "" {
			tx.Rollback()
			c.SetSubResult(model.EquitySubResult{
				Reason:  "不在限购白名单中",
				Step:    "fail",
				OrderNo: req.OrderNo,
			})
			return
		}
		limitNum = user.LimitNum
	} else {
		limitNum = activityInfo.LimitBuy
	}
	// 判断购买数量
	alreadyBuyNum, err := g.DB().Model("equity_orders").
		Where("user_id = ?", req.UserId).
		Where("activity_id = ?", req.Id).
		Count()
	if err != nil {
		tx.Rollback()
		c.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	if alreadyBuyNum >= limitNum {
		tx.Rollback()
		c.SetSubResult(model.EquitySubResult{
			Reason:  "超过限定购买数量",
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	// 生成订单
	order := model.EquityOrder{
		PublisherId: req.PublisherId,
		OrderNo:     req.OrderNo,
		Num:         req.Num,
		RealFee:     req.Num * activityInfo.Price,
		ActivityId:  activityInfo.Id,
		UserId:      req.UserId,
		Status:      1,
		Price:       activityInfo.Price,
		PayExpireAt: gtime.Now().Add(time.Minute * 10),
	}
	err = EquityOrder.Create(tx, order)
	if err != nil {
		tx.Rollback()
		c.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	return
}

func (c *equity) SetSubResult(in model.EquitySubResult) {
	_, err := g.Redis().Do("SET", fmt.Sprintf(model.SubSetEquityResultKey, in.OrderNo), gconv.String(in), "ex", 3600)
	if err != nil {
		g.Log().Errorf("EquityBuy err:%v", err)
		return
	}
	return
}

func (c *equity) GetValidDetail(id int) (ret model.EquityActivity, err error) {
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
