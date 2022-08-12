package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"meta_launchpad/model"
)

type adminSubscribeActivity struct {
}

var AdminSubscribeActivity = new(adminSubscribeActivity)

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
