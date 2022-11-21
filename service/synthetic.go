package service

import (
	"context"
	"encoding/json"
	"fmt"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"strings"
	"time"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
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
	//err = s.Check(ret.SyntheticActivity)
	//if err != nil {
	//	return
	//}
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
		ret.ConditionArr[k].Name = tamplateInfos[v.AppId+v.TemplateId].AssetName
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
	now := time.Now()
	ret.AssetPic = assetDetail.AssetPic
	ret.AssetName = assetDetail.AssetName
	ret.ChainName = publisherInfo.ChainName
	ret.ChainAddr = publisherInfo.ChainAddr
	ret.ChainType = publisherInfo.ChainType
	ret.AssetTotal = assetDetail.Total
	ret.AssetCreateAt = assetDetail.CreateTime
	ret.AssetDetailImg = templateInfo.DetailImg

	if now.Unix() >= ret.SyntheticActivity.StartTime.Unix() && now.Unix() <= ret.SyntheticActivity.EndTime.Unix() {
		ret.StatusTxt = "进行中"
	} else {
		if now.Unix() <= ret.SyntheticActivity.StartTime.Unix() {
			ret.StatusTxt = "未开始"
		} else {
			ret.StatusTxt = "已结束"
		}
	}

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

func (s *synthetic) UpdateRemainNum(id int, num int) (err error) {
	_, err = g.DB().Exec("UPDATE synthetic_activity SET remain_num = remain_num - ? WHERE id = ?", num, id)
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

func (s *synthetic) List(publisherId string, pageNum, pageSize int, startTimeBegin, startTimeEnd, endTimeBegin, endTimeEnd, status, searchVal string, orderBy string, open int, client bool) (ret model.SyntheticActivityList, err error) {
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

func (s *synthetic) SetResult(in model.SyntheticRet) {
	key := fmt.Sprintf("synthetic:%s", in.OrderNo)
	g.Redis().Do("SET", key, gconv.String(in), "ex", 7200)
	return
}

func (s *synthetic) GetResult(orderNo string) (ret model.SyntheticRet) {
	key := fmt.Sprintf("synthetic:%s", orderNo)
	gv, _ := g.Redis().DoVar("GET", key)
	if gv != nil && !gv.IsEmpty() {
		gv.Scan(&ret)
	}
	return
}

func (s *synthetic) CreateRecord(in model.SyntheticRecord) (err error) {
	_, err = g.DB().Model("synthetic_record").Insert(&in)
	return
}

func (s *synthetic) GetRecordList(pageNum, pageSize int, publisherId string, userId string, orderNo string, aid int) (ret model.SyntheticRecordList, err error) {
	m := g.DB().Model("synthetic_record")
	if aid != 0 {
		m = m.Where("aid = ?", aid)
	}
	if publisherId != "" {
		m = m.Where("publisher_id = ?", publisherId)
	}
	if userId != "" {
		m = m.Where("user_id = ?", userId)
	}
	if orderNo != "" {
		m = m.Where("order_no = ?", orderNo)
	}
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

func (s *synthetic) Synthetic(in model.SyntheticReq) {
	ainfo, err := s.ClientDetail(in.Aid)
	if err != nil {
		s.SetResult(model.SyntheticRet{
			Step:    "fail",
			OrderNo: in.OrderNo,
			Reason:  err.Error(),
		})
		return
	}
	if ainfo.RemainNum <= 0 {
		s.SetResult(model.SyntheticRet{
			Step:    "fail",
			OrderNo: in.OrderNo,
			Reason:  "合成条件不符合",
		})
		return
	}
	inMap := make(map[string]int)
	cMap := make(map[string]int)
	knIds := make([]int, 0)
	for _, vv := range ainfo.ConditionArr {
		onwer, _ := provider.KnapsackService.GetListByTemplate(in.UserId, vv.AppId, vv.TemplateId)
		for _, v := range onwer.List {
			for _, j := range in.ConditionArr {
				if v.AppId == j.AppId && v.TokenId == j.TokenId {
					inMap[v.Metadata.TemplateId]++
					knIds = append(knIds, v.Id)
				}
			}
		}
		cMap[vv.TemplateId] = vv.Num
	}
	//fmt.Println(cMap, inMap)
	g.Log().Infof("合成：cMap %+v\n", cMap)
	g.Log().Infof("合成：inMap %+v\n", inMap)
	for templateId, num := range cMap {
		if inMap[templateId] != num {
			s.SetResult(model.SyntheticRet{
				Step:    "fail",
				OrderNo: in.OrderNo,
				Reason:  "合成条件不符合",
			})
			return
		}
	}
	err = provider.KnapsackService.DeleteByIds(knIds, "SAAS_HC", "SAAS合成")
	if err != nil {
		s.SetResult(model.SyntheticRet{
			Step:    "fail",
			OrderNo: in.OrderNo,
			Reason:  err.Error(),
		})
		return
	}

	assets, e := provider.Asset.GetCanUsedAssetsByTemplate(&map[string]interface{}{
		"appId":      ainfo.AppId,
		"templateId": ainfo.TemplateId,
		"num":        ainfo.OutNum,
	})
	if e != nil && !strings.Contains(err.Error(), "timeout") {
		s.SetResult(model.SyntheticRet{
			Step:    "fail",
			OrderNo: in.OrderNo,
			Reason:  err.Error(),
		})
		return
	}
	if len(assets) == 0 {
		s.SetResult(model.SyntheticRet{
			Step:    "fail",
			OrderNo: in.OrderNo,
			Reason:  "库存不足",
		})
		return
	}
	list := make([]map[string]interface{}, 0)
	for _, v := range assets {
		item := make(map[string]interface{})
		item["userId"] = in.UserId
		item["appId"] = ainfo.AppId
		item["tokenId"] = v.TokenId
		list = append(list, item)
	}
	params := &map[string]interface{}{
		"list": list,
		"opt":  map[string]interface{}{"optUserId": in.UserId, "optType": "SAAS_HC_" + in.OrderNo, "optRemark": "SAAS合成发放资产"},
	}
	_, err = utils.SendJsonRpc(context.Background(), "knapsack", "AssetKnapsack.Add", params)
	if err != nil && !strings.Contains(err.Error(), "timeout") {
		s.SetResult(model.SyntheticRet{
			Step:    "fail",
			OrderNo: in.OrderNo,
			Reason:  err.Error(),
		})
		return
	}

	tplInfo, _ := provider.Asset.GetMateDataByTpls(&map[string]interface{}{
		"appIds":      []string{ainfo.AppId},
		"templateIds": []string{ainfo.TemplateId},
	})
	err = s.UpdateRemainNum(in.Aid, ainfo.OutNum)
	if err != nil {
		s.SetResult(model.SyntheticRet{
			Step:    "fail",
			OrderNo: in.OrderNo,
			Reason:  err.Error(),
		})
		return
	}
	var record model.SyntheticRecord
	record.OrderNo = in.OrderNo
	record.AssetName = tplInfo[ainfo.AppId+ainfo.TemplateId].AssetName
	record.AssetIcon = tplInfo[ainfo.AppId+ainfo.TemplateId].Icon
	record.AssetPic = tplInfo[ainfo.AppId+ainfo.TemplateId].AssetPic
	record.Aid = in.Aid
	record.InData = gconv.String(in.ConditionArr)
	record.OutData = gconv.String(assets)
	record.UserId = in.UserId
	record.PublisherId = in.PublisherId
	err = s.CreateRecord(record)
	if err != nil {
		s.SetResult(model.SyntheticRet{
			Step:    "fail",
			OrderNo: in.OrderNo,
			Reason:  err.Error(),
		})
		return
	}

	s.SetResult(model.SyntheticRet{
		Step:    "success",
		OrderNo: in.OrderNo,
		Reason:  "合成成功",
	})
}
