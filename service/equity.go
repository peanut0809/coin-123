package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"meta_launchpad/provider"
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
	timeNow := time.Now()
	for _, v := range equity {
		res.List = append(res.List, v)
		if v.ActivityStartTime.Unix() > timeNow.Unix() {
			v.ActivityStatus = 0
			v.ActivityStatusTxt = "未开始"
		}
		if timeNow.Unix() > v.ActivityStartTime.Unix() && timeNow.Unix() < v.ActivityEndTime.Unix() {
			v.ActivityStatus = 1
			v.ActivityStatusTxt = "进行中"
		}
		if timeNow.Unix() > v.ActivityEndTime.Unix() {
			v.ActivityStatus = 2
			v.ActivityStatusTxt = "已结束"
		}
	}
	return
}

// 活动详情
func (c *equity) Info(activityId int) (res model.EquityActivityFull, err error) {
	m := g.DB().Model("equity_activity")
	m.Where("id = ?", activityId)
	err = m.Scan(&res)
	res.LastSec = res.ActivityStartTime.Unix() - time.Now().Unix()
	res.PriceYuan = fmt.Sprintf("%.2f", float64(res.Price)/100)
	if res.LastSec < 0 {
		res.LastSec = 0
	}
	timeNow := time.Now()
	if res.ActivityStartTime.Unix() > timeNow.Unix() {
		res.ActivityStatus = 0
		res.ActivityStatusTxt = "未开始"
	}
	if timeNow.Unix() > res.ActivityStartTime.Unix() && timeNow.Unix() < res.ActivityEndTime.Unix() {
		res.ActivityStatus = 1
		res.ActivityStatusTxt = "进行中"
	}
	if timeNow.Unix() > res.ActivityEndTime.Unix() {
		res.ActivityStatus = 2
		res.ActivityStatusTxt = "已结束"
	}
	//获取应用信息
	appInfo, _ := provider.Developer.GetAppInfo(res.AppId)
	res.AssetCateString = appInfo.Data.CnName
	//获取资产分类
	templateInfo, _ := provider.Developer.GetAssetsTemplate(res.AppId, res.TemplateId)
	for _, v := range templateInfo.CateList {
		res.AssetCateString += fmt.Sprintf("-%s", v.CnName)
	}
	assetDetail, e := provider.Asset.GetMateDataByAm(&map[string]interface{}{
		"appId":      res.AppId,
		"templateId": res.TemplateId,
	})
	if e != nil {
		g.Log().Errorf("资产详情错误：%v", e)
		err = fmt.Errorf("获取资产信息失败")
		return
	}
	publisherInfo, e := provider.Developer.GetPublishInfo(res.PublisherId)
	if e != nil {
		g.Log().Errorf("发行商异常：%v，%s", e, res.PublisherId)
		err = fmt.Errorf("获取发行商失败")
		return
	}
	res.AssetPic = assetDetail.AssetPic
	res.ChainName = publisherInfo.ChainName
	res.ChainAddr = publisherInfo.ChainAddr
	res.ChainType = publisherInfo.ChainType
	res.AssetTotal = assetDetail.Total
	res.AssetCreateAt = assetDetail.CreateTime
	res.AssetDetailImg = templateInfo.DetailImg
	if templateInfo.CopyrightOpen == 1 {
		res.CopyrightInfo = templateInfo.CopyrightInfoJson
	}
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
		EquityOrder.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	// 创建订单
	err := EquityOrder.Create(&req, activityInfo)
	if err != nil {
		EquityOrder.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	return
}

// 获取活动详情
func (c *equity) GetValidDetail(id int) (ret *model.EquityActivity, err error) {
	now := time.Now()
	err = g.DB().Model("equity_activity").Where("id = ?", id).Scan(&ret)
	if err != nil {
		return
	}
	if ret == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	if now.Unix() > ret.ActivityStartTime.Unix() && now.Unix() < ret.ActivityEndTime.Unix() {
		ret.Status = model.EquityActivityStatusIng
	} else {
		if now.Unix() < ret.ActivityStartTime.Unix() {
			ret.Status = model.EquityActivityStatusWait
			err = fmt.Errorf("活动暂未开始")
			//ret.LastSec = as.ActivityStartTime.Unix() - now.Unix()
		} else {
			err = fmt.Errorf("活动已结束")
			ret.Status = model.EquityActivityStatusEnd
		}
	}
	return ret, err
}
