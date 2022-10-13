package service

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"time"
)

type synthetic struct {
}

var Synthetic = new(synthetic)

func (s *synthetic) Create(in model.SyntheticActivity) (err error) {
	in.Open = 1
	m := g.DB().Model("synthetic_activity")
	_, err = m.Insert(&in)
	if err != nil {
		return
	}
	return
}

func (s *synthetic) Detail(id int) (ret model.SyntheticActivity, err error) {
	m := g.DB().Model("synthetic_activity")
	err = m.Where("id", id).Scan(&ret)
	if err != nil {
		return
	}
	return
}

func (s *synthetic) Check(in model.SyntheticActivity) (err error) {
	now := time.Now()
	if !(now.Unix() >= in.StartTime.Unix() && now.Unix() <= in.EndTime.Unix()) {
		err = fmt.Errorf("活动不在进行中")
		return
	}
	if in.Open != 1 {
		err = fmt.Errorf("活动已关闭")
		return
	}
	return
}

func (s *synthetic) ClientDetail(id int) (ret model.SyntheticActivityDetail, err error) {
	ret.SyntheticActivity, err = s.Detail(id)
	if err != nil {
		return
	}
	err = s.Check(ret.SyntheticActivity)
	if err != nil {
		return
	}
	if ret.Condition != nil {
		json.Unmarshal([]byte(*ret.Condition), &ret.ConditionArr)
	}
	appIds := make([]string, 0)
	templateIds := make([]string, 0)
	for _, v := range ret.ConditionArr {
		appIds = append(appIds, v.AppId)
		templateIds = append(templateIds, v.TemplateId)
	}
	tamplateInfos, _ := provider.Asset.GetMateDataByTpls(&map[string]interface{}{
		"appIds":      appIds,
		"templateIds": templateIds,
	})
	for k, v := range ret.ConditionArr {
		ret.ConditionArr[k].Cover = tamplateInfos[v.AppId+v.TemplateId].Icon
	}

	//获取应用信息
	appInfo, _ := provider.Developer.GetAppInfo(ret.SyntheticActivity.AppId)
	ret.AssetCateString = appInfo.Data.CnName
	//获取资产分类
	templateInfo, _ := provider.Developer.GetAssetsTemplate(ret.SyntheticActivity.AppId, ret.SyntheticActivity.TemplateId)
	for _, v := range templateInfo.CateList {
		ret.AssetCateString += fmt.Sprintf("-%s", v.CnName)
	}
	assetDetail, e := provider.Asset.GetMateDataByAm(&map[string]interface{}{
		"appId":      ret.SyntheticActivity.AppId,
		"templateId": ret.SyntheticActivity.TemplateId,
	})
	if e != nil {
		g.Log().Errorf("资产详情错误：%v", e)
		err = fmt.Errorf("获取资产信息失败")
		return
	}
	publisherInfo, e := provider.Developer.GetPublishInfo(ret.SyntheticActivity.PublisherId)
	if e != nil {
		g.Log().Errorf("发行商异常：%v，%s", e, ret.SyntheticActivity.PublisherId)
		err = fmt.Errorf("获取发行商失败")
		return
	}
	ret.AssetPic = assetDetail.AssetPic
	ret.AssetName = assetDetail.AssetName
	ret.ChainName = publisherInfo.ChainName
	ret.ChainAddr = publisherInfo.ChainAddr
	ret.ChainType = publisherInfo.ChainType
	ret.AssetTotal = assetDetail.Total
	ret.AssetCreateAt = assetDetail.CreateTime
	ret.AssetDetailImg = templateInfo.DetailImg
	return
}

func (s *synthetic) Open(id int, open int) (ret model.SyntheticActivity, err error) {
	m := g.DB().Model("synthetic_activity")
	_, err = m.Where("id", id).Data(g.Map{
		"open": open,
	}).Update()
	if err != nil {
		return
	}
	return
}

func (s *synthetic) Delete(id int) (ret model.SyntheticActivity, err error) {
	rcount, _ := g.DB().Model("synthetic_record").Where("aid = ?", id).Count()
	if rcount != 0 {
		err = fmt.Errorf("活动已被用户合成，不能删除")
		return
	}
	m := g.DB().Model("synthetic_activity")
	_, err = m.Where("id", id).Delete()
	if err != nil {
		return
	}
	return
}

func (s *synthetic) Update(in model.SyntheticActivity) (err error) {
	//有人合成就不能修改
	rcount, _ := g.DB().Model("synthetic_record").Where("aid = ?", in.Id).Count()
	if rcount != 0 {
		err = fmt.Errorf("活动已被用户合成，不能修改")
		return
	}
	m := g.DB().Model("synthetic_activity")
	_, err = m.Data(g.Map{
		"name":        in.Name,
		"app_id":      in.AppId,
		"asset_type":  in.AssetType,
		"template_id": in.TemplateId,
		"sum":         in.Sum,
		"remain_num":  in.RemainNum,
		"out_num":     in.OutNum,
		"cover":       in.Cover,
		"rule":        in.Rule,
		"start_time":  in.StartTime,
		"end_time":    in.EndTime,
		"condition":   in.Condition,
	}).Where("id", in.Id).Update()
	if err != nil {
		return
	}
	return
}

func (s *synthetic) List(publisherId string, pageNum, pageSize int, startTimeBegin, startTimeEnd, endTimeBegin, endTimeEnd, status, searchVal string, orderBy string, open int) (ret model.SyntheticActivityList, err error) {
	m := g.DB().Model("synthetic_activity").Where("publisher_id", publisherId)
	if startTimeBegin != "" && startTimeEnd != "" {
		m = m.Where("start_time >= ? AND start_time <= ?", startTimeBegin, startTimeEnd)
	}
	if endTimeBegin != "" && endTimeEnd != "" {
		m = m.Where("end_time >= ? AND end_time <= ?", endTimeBegin, endTimeEnd)
	}
	if open != 0 {
		m = m.Where("open", open)
	}
	now := time.Now()
	if status != "" {
		if status == "ing" {
			m = m.Where("? >= start_time AND ? <= end_time", now, now)
		}
		if status == "waitStart" {
			m = m.Where("? <= start_time", now)
		}
		if status == "end" {
			m = m.Where("? >= end_time", now)
		}
	}
	if searchVal != "" {
		m = m.Where("(name LIKE ?) OR (id = ?)", "%"+searchVal+"%", searchVal)
	}
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	if orderBy != "" {
		m = m.Order(orderBy)
	}
	err = m.Page(pageNum, pageSize).Scan(&ret.List)
	if err != nil {
		return
	}
	for k, v := range ret.List {
		if now.Unix() >= v.StartTime.Unix() && now.Unix() <= v.EndTime.Unix() {
			ret.List[k].StatusTxt = "进行中"
		} else {
			if now.Unix() <= v.StartTime.Unix() {
				ret.List[k].StatusTxt = "未开始"
			} else {
				ret.List[k].StatusTxt = "已结束"
			}
		}
	}
	return
}
