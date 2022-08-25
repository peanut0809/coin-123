package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"meta_launchpad/model"
	"meta_launchpad/provider"
)

type adminSubscribeActivity struct {
}

var AdminSubscribeActivity = new(adminSubscribeActivity)

func (s *adminSubscribeActivity) ListByPage(activityType int, publisherId string, pageNum int, createStartTime, createEndTime, activityStartTimeA, activityStartTimeB, status, activityEndTimeA, activityEndTimeB, searchVal string, pageSize int) (ret model.AdminListByPage, err error) {
	m := g.DB().Model("subscribe_activity").Where("publisher_id = ? AND activity_type = ?", publisherId, activityType)
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
	list := make([]model.SubscribeActivity, 0)
	err = m.Order("id DESC").Page(pageNum, pageSize).Scan(&list)
	if err != nil {
		return
	}
	for _, v := range list {
		item := model.AdminSubscribeActivityFull{
			SubscribeActivity: v,
			Status:            "",
			StatusTxt:         "",
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

func (s *adminSubscribeActivity) Create(in model.SubscribeActivity, cons []model.SubscribeCondition) (err error) {
	var tx *gdb.TX
	tx, err = g.DB().Begin()
	if err != nil {
		return
	}
	var sqlRet sql.Result
	m := tx.Model("subscribe_activity")
	in.TicketInfo = `[
		{
			"use": false,
			"type": "money",
			"unitNum": 1
		},
		{
			"use": true,
			"type": "crystal",
			"unitNum": 0
		},
		{
			"use": false,
			"type": "month_ticket",
			"unitNum": 1
		}
	]`
	in.StartTime = gtime.Now()
	in.RemainNum = in.SumNum
	sqlRet, err = m.Insert(&in)
	if err != nil {
		tx.Rollback()
		return
	}
	aid, e := sqlRet.LastInsertId()
	if e != nil {
		err = e
		tx.Rollback()
		return
	}
	in.Id = int(aid)
	if in.ActivityType == 2 { //普通购
		err = tx.Commit()
		return
	}
	//优先购
	if len(cons) != 0 {
		m := tx.Model("subscribe_condition")
		for _, v := range cons {
			if v.AppId == "" {
				tx.Rollback()
				err = fmt.Errorf("appId不能为空")
				return
			}
			if v.AssetType == "" {
				tx.Rollback()
				err = fmt.Errorf("AssetType不能为空")
				return
			}
			if v.TemplateId == "" {
				tx.Rollback()
				err = fmt.Errorf("TemplateId不能为空")
				return
			}
			if v.BuyNum <= 0 {
				tx.Rollback()
				err = fmt.Errorf("购买数量参数错误")
				return
			}
			if v.MetaDataRule != "" {
				vJson := make(map[string]string)
				err = json.Unmarshal([]byte(v.MetaDataRule), &vJson)
				if err != nil {
					tx.Rollback()
					err = fmt.Errorf("MetaDataRule参数不合法")
					return
				}
				if len(vJson) == 0 {
					tx.Rollback()
					err = fmt.Errorf("MetaDataRule参数不合法")
					return
				}
			} else {
				m = m.FieldsEx("meta_data_rule")
			}
			v.PublisherId = in.PublisherId
			v.Aid = in.Id
			_, err = m.Insert(&v)
			if err != nil {
				tx.Rollback()
				return
			}
		}
		err = tx.Commit()
		return
	} else {
		err = tx.Commit()
		return
	}
}

func (s *adminSubscribeActivity) Update(in model.SubscribeActivity, cons []model.SubscribeCondition) (err error) {
	count := SubscribeRecord.GetSubscribeRecordCount(in.Id)
	if count > 0 {
		err = fmt.Errorf("活动已被认购，不能修改")
		return
	}
	var tx *gdb.TX
	tx, err = g.DB().Begin()
	if err != nil {
		return
	}
	m := tx.Model("subscribe_activity")
	in.TicketInfo = `[
		{
			"use": false,
			"type": "money",
			"unitNum": 1
		},
		{
			"use": true,
			"type": "crystal",
			"unitNum": 0
		},
		{
			"use": false,
			"type": "month_ticket",
			"unitNum": 1
		}
	]`
	in.StartTime = gtime.Now()
	in.RemainNum = in.SumNum
	updateMap := g.Map{
		"name":                in.Name,
		"activity_start_time": in.ActivityStartTime,
		"activity_end_time":   in.ActivityEndTime,
		"price":               in.Price,
		"activity_intro":      in.ActivityIntro,
		"cover_img_url":       in.CoverImgUrl,
		"app_id":              in.AppId,
		"asset_type":          in.AssetType,
		"asset_type2":         in.AssetType2,
		"asset_type3":         in.AssetType3,
		"template_id":         in.TemplateId,
		"sum_num":             in.SumNum,
		"remain_num":          in.SumNum,
		"open_award_time":     in.OpenAwardTime,
		"pay_end_time":        in.PayEndTime,
		"nfr_sec":             in.NfrSec,
	}
	if in.ActivityType == 2 {
		updateMap["general_buy_num"] = in.GeneralBuyNum
	}
	if in.ActivityType == 1 {
		updateMap["award_method"] = in.AwardMethod
	}
	_, err = m.Data(updateMap).Where("id", in.Id).Update()
	if err != nil {
		tx.Rollback()
		return
	}

	if in.ActivityType == 2 { //普通购
		err = tx.Commit()
		return
	}
	//优先购

	//删除所有条件
	_, err = tx.Model("subscribe_condition").Where("aid", in.Id).Delete()
	if err != nil {
		tx.Rollback()
		return
	}
	if len(cons) != 0 {
		m := tx.Model("subscribe_condition")
		for _, v := range cons {
			if v.AppId == "" {
				tx.Rollback()
				err = fmt.Errorf("appId不能为空")
				return
			}
			if v.AssetType == "" {
				tx.Rollback()
				err = fmt.Errorf("AssetType不能为空")
				return
			}
			if v.TemplateId == "" {
				tx.Rollback()
				err = fmt.Errorf("TemplateId不能为空")
				return
			}
			if v.BuyNum <= 0 {
				tx.Rollback()
				err = fmt.Errorf("购买数量参数错误")
				return
			}
			if v.MetaDataRule != "" {
				vJson := make(map[string]string)
				err = json.Unmarshal([]byte(v.MetaDataRule), &vJson)
				if err != nil {
					tx.Rollback()
					err = fmt.Errorf("MetaDataRule参数不合法")
					return
				}
				if len(vJson) == 0 {
					tx.Rollback()
					err = fmt.Errorf("MetaDataRule参数不合法")
					return
				}
			} else {
				m = m.FieldsEx("meta_data_rule")
			}
			v.PublisherId = in.PublisherId
			v.Aid = in.Id
			_, err = m.Insert(&v)
			if err != nil {
				tx.Rollback()
				return
			}
		}
		err = tx.Commit()
		return
	} else {
		err = tx.Commit()
		return
	}
}

func (s *adminSubscribeActivity) Detail(publisherId string, id int) (ret model.AdminSubscribeActivityDetail, err error) {
	var activiInfo *model.SubscribeActivity
	err = g.DB().Model("subscribe_activity").Where("id = ? AND publisher_id = ?", id, publisherId).Scan(&activiInfo)
	if err != nil {
		return
	}
	if activiInfo == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	ret.SubscribeActivity = *activiInfo
	ret.PriceYuan = fmt.Sprintf("%.2f", float64(ret.Price)/100)
	if ret.ActivityType == 1 { //如果是优先购，查出条件
		err = g.DB().Model("subscribe_condition").Where("aid", ret.Id).Scan(&ret.Cons)
		if err != nil {
			return
		}
		if len(ret.Cons) == 0 {
			ret.Cons = make([]model.SubscribeCondition, 0)
		}
	}
	return
}

func (s *adminSubscribeActivity) Delete(publisherId string, id int) (err error) {
	err = g.DB().Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) (err error) {
		_, err = tx.Model("subscribe_activity").Where("id = ? AND publisher_id = ?", id, publisherId).Delete()
		if err != nil {
			return
		}
		_, err = tx.Model("subscribe_condition").Where("aid", id).Delete()
		if err != nil {
			return
		}
		return
	})
	return
}

func (s *adminSubscribeActivity) Disable(id int, disable int, publisherId string) (err error) {
	var r sql.Result
	r, err = g.DB().Exec("UPDATE subscribe_activity SET disable = ? WHERE id = ? AND publisher_id = ?", disable, id, publisherId)
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

func (s *adminSubscribeActivity) GetSubRecords(activityType int, pageNum int, publisherId string, createdAtStart, createdAtEnd string, priceMin int, priceMax int, award int, payStatus int, searchVal string) (ret model.AdminSubscribeRecordByPage, err error) {
	m := g.DB().Model("subscribe_records").Where("publisher_id = ? AND activity_type = ?", publisherId, activityType)
	if createdAtStart != "" && createdAtEnd != "" {
		m = m.Where("created_at >= ? AND created_at <= ?", createdAtStart, createdAtEnd)
	}
	if priceMin != 0 && priceMax != 0 {
		m = m.Where("sum_price >= ? AND sum_price <= ?", priceMin, priceMax)
	}
	if award != -1 {
		m = m.Where("award = ?", award)
	}
	if payStatus != -1 {
		m = m.Where("pay_status = ?", payStatus)
	}
	if searchVal != "" {
		m = m.Where("(aid = ? OR order_no = ? OR pay_order_no = ? OR name LIKE ?)", searchVal, searchVal, searchVal, "%"+searchVal+"%")
	}
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	var list []model.SubscribeRecord
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
		item := model.AdminSubscribeRecordFull{
			SubscribeRecord: v,
			OrderStatusTxt:  "",
			UserPhone:       userInfoMap[v.UserId].Phone,
			UserName:        userInfoMap[v.UserId].Nickname,
			SumPriceYuan:    fmt.Sprintf("%.2f", float64(v.SumPrice)/100),
		}
		if v.PayStatus == 0 {
			item.OrderStatusTxt = "未支付"
		}
		if v.PayStatus == 1 {
			item.OrderStatusTxt = "已支付"
		}
		if v.PayStatus == 2 {
			item.OrderStatusTxt = "已超时"
		}
		ret.List = append(ret.List, item)
	}
	return
}
