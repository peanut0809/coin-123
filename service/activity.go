package service

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"context"
	"encoding/json"
	"fmt"
	"meta_launchpad/model"
	"meta_launchpad/provider"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type activity struct {
}

var Activity = new(activity)

func (s *activity) GetByIds(ids []int) (ret []model.Activity) {
	_ = g.DB().Model("activity").Where("id in (?)", ids).Scan(&ret)
	return
}

func (s *activity) GetCreatorRank(rankValue int, pageNum int, pageSize int, publisherId string, searchVal string) (ret map[string]interface{}, err error) {
	ret = make(map[string]interface{})
	m := g.DB().Model("subscribe_activity").Where("publisher_id = ?", publisherId)
	if rankValue != 0 {
		m = m.Where("sum_num = ?", rankValue)
	}
	if searchVal != "" {
		m = m.Where("(creator_name LIKE ? OR creator_no = ?)", "%"+searchVal+"%", searchVal)
	}
	ret["total"], err = m.Count()
	if err != nil {
		return
	}
	var list []model.SubscribeActivity
	err = m.Order("sum_num DESC,price DESC").Page(pageNum, pageSize).Scan(&list)
	if err != nil {
		return
	}
	ret["list"] = list
	return
}

func (s *activity) List(activityIds []int, pageNum int, pageSize int, startTime, endTime string, activityType int, status, searchVal, publisherId string, disable int) (ret model.AdminActivityList, err error) {
	m := g.DB().Model("activity")
	if disable != -1 {
		m = m.Where("disable = ?", disable)
	}
	if publisherId != "" {
		m = m.Where("publisher_id = ?", publisherId)
	}
	if startTime != "" {
		m = m.Where("start_time >= ?", startTime)
	}
	if endTime != "" {
		m = m.Where("end_time <= ?", endTime)
	}
	if len(activityIds) != 0 {
		m = m.Where("id IN (?)", activityIds)
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
	err = m.Order("id DESC").Page(pageNum, pageSize).Scan(&as)
	if err != nil {
		return
	}
	subIds := make([]int, 0)
	secKillIds := make([]int, 0)
	equityIds := make([]int, 0)
	for _, v := range as {
		if v.ActivityType == 3 {
			secKillIds = append(secKillIds, v.ActivityId)
		} else if v.ActivityType == 4 {
			equityIds = append(equityIds, v.ActivityId)
		} else {
			subIds = append(subIds, v.ActivityId)
		}
	}
	var (
		subAcMap     map[int]model.SubscribeActivity
		secKillAcMap map[int]model.SeckillActivity
		equityAcMap  map[int]model.EquityActivity
	)
	if len(subIds) != 0 {
		subAcMap = SubscribeActivity.GetByIds(subIds)
	}
	if len(secKillIds) != 0 {
		secKillAcMap = SeckillActivity.GetByIds(secKillIds)
	}
	if len(equityIds) != 0 {
		equityAcMap = Equity.GetByIds(equityIds)
	}
	type developers struct {
		Id             int    `json:"id"`
		RelationUserId string `json:"relationUserId"`
		Name           string `json:"name"`
		LogoUrl        string `json:"logoUrl"`
	}
	var developerDetails []developers
	developer := make(map[int]developers)
	var developerIds []int
	for _, i := range as {
		developerIds = append(developerIds, subAcMap[i.ActivityId].CreatorId)
	}
	rpc, err := utils.SendJsonRpc(context.Background(), "developer", "Developer.DeveloperList", g.Map{
		"ids": developerIds,
	})
	if err != nil {
		return
	}
	mar, _ := json.Marshal(rpc)
	_ = json.Unmarshal(mar, &developerDetails)
	for _, i := range developerDetails {
		developer[i.Id] = i
	}

	publisherIds := make([]string, 0)
	for _, v := range as {
		publisherIds = append(publisherIds, v.PublisherId)
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
		} else if v.ActivityType == 4 {
			item.SumNum = equityAcMap[v.ActivityId].Number
			item.Price = fmt.Sprintf("%.2f", float64(equityAcMap[v.ActivityId].Price)/100)
			item.Cover = equityAcMap[v.ActivityId].CoverImgUrl
			item.ActivityTypeString = "权益购"
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
			if publisherId == "MCN" {
				item.PublisherName = developer[subAcMap[v.ActivityId].CreatorId].Name
				item.PublisherIcon = developer[subAcMap[v.ActivityId].CreatorId].LogoUrl
				item.PublisherUserId = developer[subAcMap[v.ActivityId].CreatorId].RelationUserId
				//item.PublisherName = subAcMap[v.ActivityId].CreatorName
				//item.PublisherIcon = subAcMap[v.ActivityId].CreatorAvatar
				//item.PublisherUserId = subAcMap[v.ActivityId].CreatorUserId
			}
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
	if publisherId == "" {
		publisherInfoMap, _ := provider.Developer.GetPublisherByIds(publisherIds)
		for i := range ret.List {
			ret.List[i].PublisherName = publisherInfoMap[ret.List[i].PublisherId].Name
			ret.List[i].PublisherIcon = publisherInfoMap[ret.List[i].PublisherId].Icon
		}
	}
	return
}
