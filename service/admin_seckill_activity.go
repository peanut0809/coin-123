package service

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"meta_launchpad/model"
	"meta_launchpad/provider"
)

type adminSecKillActivity struct {
}

var AdminSecKillActivity = new(adminSecKillActivity)

func (s *adminSecKillActivity) Create(in model.SeckillActivity) (err error) {
	_, err = g.DB().Model("seckill_activity").Insert(&in)
	if err != nil {
		return
	}
	return
}

func (s *adminSecKillActivity) GetOrders(pageNum int, publisherId string, createdAtStart, createdAtEnd string, priceMin int, priceMax int, payStatus int, searchVal string) (ret model.AdminSeckillOrderByPage, err error) {
	m := g.DB().Model("seckill_orders").Where("publisher_id = ?", publisherId)
	if createdAtStart != "" && createdAtEnd != "" {
		m = m.Where("created_at >= ? AND created_at <= ?", createdAtStart, createdAtEnd)
	}
	if priceMin != 0 && priceMax != 0 {
		m = m.Where("real_fee >= ? AND real_fee <= ?", priceMin, priceMax)
	}
	if payStatus != 0 {
		m = m.Where("pay_status = ?", payStatus)
	}
	if searchVal != "" {
		m = m.Where("(aid = ? OR order_no = ? OR name LIKE ?)", searchVal, searchVal, "%"+searchVal+"%")
	}
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	var list []model.SeckillOrder
	err = m.Order("id DESC").Page(pageNum, 20).Scan(&list)
	if err != nil {
		return
	}
	userIds := make([]string, 0)
	for _, v := range list {
		userIds = append(userIds, v.UserId)
	}
	_, userInfoMap, _ := provider.User.GetUserInfo(userIds)
	for _, v := range list {
		item := model.AdminSeckillOrderFull{
			SeckillOrder: v,
			StatusTxt:    "",
			UserName:     userInfoMap[v.UserId].Nickname,
			UserPhone:    userInfoMap[v.UserId].Phone,
			RealFeeYuan:  fmt.Sprintf("%.2f", float64(v.RealFee)/100),
		}
		//1.待支付；2.已支付；3.已超时；4.已取消
		if v.Status == 1 {
			item.StatusTxt = "待支付"
		}
		if v.Status == 1 {
			item.StatusTxt = "已支付"
		}
		if v.Status == 3 {
			item.StatusTxt = "已超时"
		}
		if v.Status == 4 {
			item.StatusTxt = "已取消"
		}
		ret.List = append(ret.List, item)
	}
	return
}

func (s *adminSecKillActivity) Update(in model.SeckillActivity) (err error) {
	updateMap := g.Map{
		"name":                in.Name,
		"app_id":              in.AppId,
		"template_id":         in.TemplateId,
		"price":               in.Price,
		"sum_num":             in.SumNum,
		"remain_num":          in.SumNum,
		"cover_img_url":       in.CoverImgUrl,
		"limit_buy":           in.LimitBuy,
		"activity_intro":      in.ActivityIntro,
		"activity_start_time": in.ActivityStartTime,
		"activity_end_time":   in.ActivityEndTime,
	}
	_, err = g.DB().Model("seckill_activity").Data(updateMap).Where("id = ? AND publisher_id = ?", in.Id, in.PublisherId).Update()
	if err != nil {
		return
	}
	return
}

func (s *adminSecKillActivity) Detail(id int, publisherId string) (ret model.CreateSeckillActivityReq, err error) {
	var activityInfo *model.SeckillActivity
	err = g.DB().Model("seckill_activity").Where("id = ? AND publisher_id = ?", id, publisherId).Scan(&activityInfo)
	if err != nil {
		return
	}
	if activityInfo == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	ret.SeckillActivity = *activityInfo
	ret.PriceYuan = fmt.Sprintf("%.2f", float64(activityInfo.Price)/100)
	return
}

func (s *adminSecKillActivity) Delete(id int, publisherId string) (err error) {
	var r sql.Result
	r, err = g.DB().Exec("DELETE FROM seckill_activity WHERE id = ? AND publisher_id = ?", id, publisherId)
	if err != nil {
		return
	}
	affectedNum, _ := r.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("删除失败")
		return
	}
	return
}

func (s *adminSecKillActivity) Disable(id int, disable int, publisherId string) (err error) {
	var r sql.Result
	r, err = g.DB().Exec("UPDATE seckill_activity SET disable = ? WHERE id = ? AND publisher_id = ?", disable, id, publisherId)
	if err != nil {
		return
	}
	affectedNum, _ := r.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("更新失败")
		return
	}
	return
}

func (s *adminSecKillActivity) List(publisherId string, pageNum int, createStartTime, createEndTime, activityStartTimeA, activityStartTimeB, status, activityEndTimeA, activityEndTimeB, searchVal string) (ret model.AdminSeckillActivityList, err error) {
	m := g.DB().Model("seckill_activity").Where("publisher_id = ?", publisherId)
	if createStartTime != "" && createEndTime != "" {
		m = m.Where("created_at >= ? and created_at <= ?", createStartTime, createEndTime)
	}
	if activityStartTimeA != "" && activityStartTimeB != "" {
		m = m.Where("activity_start_time >= ? and activity_start_time <= ?", activityStartTimeA, activityStartTimeB)
	}
	if activityEndTimeA != "" && activityEndTimeB != "" {
		m = m.Where("activity_end_time >= ? and activity_end_time <= ?", activityEndTimeA, activityEndTimeB)
	}
	n := gtime.Now()
	if status == "0" { //未开始
		m = m.Where("activity_start_time > ?", n)
	}
	if status == "1" { //进行中
		m = m.Where("? > activity_start_time AND ? < activity_end_time", n, n)
	}
	if status == "2" { //已结束
		m = m.Where("? > activity_end_time", n)
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
	var as []model.SeckillActivity
	err = m.Order("id DESC").Page(pageNum, 20).Scan(&as)
	if err != nil {
		return
	}
	for _, v := range as {
		item := model.AdminSeckillActivityFull{
			SeckillActivity: v,
			PriceYuan:       "",
			StatusTxt:       "",
		}
		if v.ActivityStartTime.Unix() > n.Unix() {
			item.Status = "0"
			item.StatusTxt = "未开始"
		}
		if n.Unix() > v.ActivityStartTime.Unix() && n.Unix() < v.ActivityEndTime.Unix() {
			item.Status = "1"
			item.StatusTxt = "进行中"
		}
		if n.Unix() > v.ActivityEndTime.Unix() {
			item.Status = "2"
			item.StatusTxt = "已结束"
		}
		item.PriceYuan = fmt.Sprintf("%.2f", float64(v.Price)/100)
		ret.List = append(ret.List, item)
	}
	return
}
