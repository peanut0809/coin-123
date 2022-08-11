package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
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
	sqlRet, err = tx.Model("subscribe_activity").Insert(&in)
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
		if err != nil {
			tx.Rollback()
			return
		}
		return
	}
	//优先购
	if len(cons) != 0 {
		for k, v := range cons {
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
			}
			cons[k].PublisherId = in.PublisherId
			cons[k].Aid = in.Id
		}
		_, err = tx.Model("subscribe_condition").Insert(&cons)
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return
		}
	} else {
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return
		}
	}
	return
}
