package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"github.com/gogf/gf/net/ghttp"
	"github.com/shopspring/decimal"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

type adminSeckillActivity struct {
	api.CommonBase
}

var AdminSeckillActivity = new(adminSeckillActivity)

func (s *adminSeckillActivity) Create(r *ghttp.Request) {
	var req model.CreateSeckillActivityReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if req.Name == "" || req.AppId == "" || req.TemplateId == "" || req.SumNum <= 0 || req.CoverImgUrl == "" || req.LimitBuy <= 0 || req.ActivityIntro == "" || req.ActivityStartTime == nil || req.ActivityEndTime == nil {
		s.FailJsonExit(r, "参数错误")
		return
	}
	decimalValue, err := decimal.NewFromString(req.PriceYuan)
	if err != nil {
		s.FailJsonExit(r, "价格参数错误")
		return
	}
	decimalValue = decimalValue.Mul(decimal.NewFromFloat(100))
	priceInt := decimalValue.IntPart()
	if priceInt <= 0 {
		s.FailJsonExit(r, "价格参数错误")
		return
	}
	if req.ActivityStartTime.Unix() > req.ActivityEndTime.Unix() {
		s.FailJsonExit(r, "活动开始时间不能大于结束时间")
		return
	}
	req.Alias = utils.RandString(6)
	req.PublisherId = s.GetPublisherId(r)
	req.Price = int(priceInt)
	req.RemainNum = req.SumNum
	if req.Id == 0 {
		err = service.AdminSecKillActivity.Create(req.SeckillActivity)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	} else {
		err = service.AdminSecKillActivity.Update(req.SeckillActivity)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	}
	s.SusJsonExit(r)
}

func (s *adminSeckillActivity) Detail(r *ghttp.Request) {
	id := r.GetQueryInt("id")
	ret, err := service.AdminSecKillActivity.Detail(id, s.GetPublisherId(r))
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *adminSeckillActivity) Disable(r *ghttp.Request) {
	id := r.GetInt("id")
	action := r.GetString("action")
	if action == "ON" {
		err := service.AdminSecKillActivity.Disable(id, 0, s.GetPublisherId(r))
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	} else if action == "OFF" {
		err := service.AdminSecKillActivity.Disable(id, 1, s.GetPublisherId(r))
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	}
	s.SusJsonExit(r)
}

func (s *adminSeckillActivity) List(r *ghttp.Request) {
	createStartTime := r.GetQueryString("createStartTime")
	createEndTime := r.GetQueryString("createEndTime")
	activityStartTimeA := r.GetQueryString("activityStartTimeA")
	activityStartTimeB := r.GetQueryString("activityStartTimeB")
	status := r.GetQueryString("status")
	activityEndTimeA := r.GetQueryString("activityEndTimeA")
	activityEndTimeB := r.GetQueryString("activityEndTimeB")
	searchVal := r.GetQueryString("searchVal")
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	publisherId := s.GetPublisherId(r)
	ret, err := service.AdminSecKillActivity.List(publisherId, pageNum, createStartTime, createEndTime, activityStartTimeA, activityStartTimeB, status, activityEndTimeA, activityEndTimeB, searchVal, pageSize)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *adminSeckillActivity) Delete(r *ghttp.Request) {
	id := r.GetInt("id")
	if service.SeckillOrder.Count(id) > 0 {
		s.FailJsonExit(r, "活动已有用户参与，不能删除")
		return
	}
	publisherId := s.GetPublisherId(r)
	err := service.AdminSecKillActivity.Delete(id, publisherId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
}

func (s *adminSeckillActivity) GetOrders(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 1)
	createdAtStart := r.GetQueryString("createdAtStart")
	createdAtEnd := r.GetQueryString("createdAtEnd")
	priceMinStr := r.GetQueryString("priceMinStr")
	priceMaxStr := r.GetQueryString("priceMaxStr")
	payStatus := r.GetQueryInt("payStatus")
	searchVal := r.GetQueryString("searchVal")
	priceMinStrValue, _ := decimal.NewFromString(priceMinStr)
	priceMinStrValue = priceMinStrValue.Mul(decimal.NewFromInt(100))
	priceMinInt := priceMinStrValue.IntPart()
	priceMaxStrValue, _ := decimal.NewFromString(priceMaxStr)
	priceMaxStrValue = priceMaxStrValue.Mul(decimal.NewFromInt(100))
	priceMaxInt := priceMaxStrValue.IntPart()
	publisherId := s.GetPublisherId(r)
	ret, err := service.AdminSecKillActivity.GetOrders(pageNum, pageSize, publisherId, createdAtStart, createdAtEnd, int(priceMinInt), int(priceMaxInt), payStatus, searchVal)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}
