package service

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"strings"
)

type drop struct {
}

var Drop = new(drop)

func (s *drop) Create(in model.DropRecord) (err error) {
	var r sql.Result
	r, err = g.DB().Model("drop_record").Insert(&in)
	if err != nil {
		return
	}
	lastId, e := r.LastInsertId()
	if e != nil {
		err = e
		return
	}
	go s.StartDrop(int(lastId))
	return
}

func (s *drop) GetDetailRecordList(pageNum, pageSize, dropId int) (ret model.DropDetailRecordList, err error) {
	m := g.DB().Model("drop_detail_record").Where("drop_id = ?", dropId)
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	err = m.Order("id DESC").Page(pageNum, pageSize).Scan(&ret.List)
	if err != nil {
		return
	}
	return
}

func (s *drop) GetRecordList(pageNum int, pageSize int, createStartTime string, createEndTime string, searchVal string) (ret model.RecordList, err error) {
	m := g.DB().Model("drop_record")
	if createStartTime != "" && createEndTime != "" {
		m = m.Where("created_at BETWEEN ? AND ?", createStartTime, createEndTime)
	}
	if searchVal != "" {
		m = m.Where("(name LIKE ? OR id = ?)", "%"+searchVal+"%", searchVal)
	}
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	rs := make([]model.DropRecord, 0)
	err = m.Order("id DESC").Page(pageNum, pageSize).Scan(&rs)
	if err != nil {
		return
	}
	appIds := make([]string, 0)
	templateIds := make([]string, 0)
	for _, v := range rs {
		appIds = append(appIds, v.AppId)
		templateIds = append(templateIds, v.TemplateId)
	}
	tamplateInfos, _ := provider.Asset.GetMateDataByTpls(&map[string]interface{}{
		"appIds":      appIds,
		"templateIds": templateIds,
	})
	for _, v := range rs {
		ret.List = append(ret.List, model.DropRecordFull{
			DropRecord: v,
			AssetName:  tamplateInfos[v.AppId+v.TemplateId].AssetName,
		})
	}
	return
}

func (s *drop) UpdateDropStatus(id int, status int) (err error) {
	_, err = g.DB().Exec("UPDATE drop_record SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return
	}
	return
}

func (s *drop) CreateDetailRecord(in *model.DropDetailRecord) (err error) {
	var r sql.Result
	r, err = g.DB().Model("drop_detail_record").Insert(in)
	if err != nil {
		return
	}
	id, _ := r.LastInsertId()
	in.Id = int(id)
	return
}

func (s *drop) UpdateDetailRecordStatus(id int, status int, tokenId string, errMsg string) (err error) {
	_, err = g.DB().Exec("UPDATE drop_detail_record SET token_id = ?,status = ?,err_msg = ? WHERE id = ?", tokenId, status, errMsg, id)
	return
}

func (s *drop) StartDrop(dropId int) (err error) {
	var dropInfo *model.DropRecord
	err = g.DB().Model("drop_record").Where("id", dropId).Scan(&dropInfo)
	if err != nil {
		return
	}
	if dropInfo == nil {
		err = fmt.Errorf("空投ID不存在：%d", dropId)
		return
	}
	phoneArr := strings.Split(dropInfo.Phones, ",")
	userMap, e := provider.User.GetUserInfoByPhone(&map[string]interface{}{
		"phoneArr": phoneArr,
	})
	if e != nil {
		err = e
		return
	}
	for _, phone := range phoneArr {
		for i := 0; i < dropInfo.Num; i++ {
			var item model.DropDetailRecord
			item.UserId = userMap[phone].UserId
			item.Phone = phone
			item.DropId = dropId
			item.TemplateId = dropInfo.TemplateId
			if item.UserId == "" {
				item.Status = 2
				item.ErrMsg = "用户不存在"
			}
			item.AppId = dropInfo.AppId
			item.NfrSec = dropInfo.NfrSec
			err = s.CreateDetailRecord(&item)
			if err != nil {
				return
			}
			if item.Status == 0 {
				err = s.PublishToMq(item)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (s *drop) PublishToMq(v interface{}) (err error) {
	queueName := "drop.asset"
	mqClient, e := client.GetQueue(client.QueueConfig{
		QueueName:  queueName,
		Exchange:   queueName,
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if e != nil {
		err = e
		return
	}
	defer mqClient.Close()
	err = mqClient.Publish(v)
	if err != nil {
		return
	}
	return
}
