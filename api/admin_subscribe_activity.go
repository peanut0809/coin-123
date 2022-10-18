package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/shopspring/decimal"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"meta_launchpad/service"
)

type adminSubscribeActivity struct {
	api.CommonBase
}

var AdminSubscribeActivity = new(adminSubscribeActivity)

func (s *adminSubscribeActivity) Create(r *ghttp.Request) {
	var req model.CreateSubscribeActivityReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if req.Name == "" || req.CoverImgUrl == "" || req.AppId == "" || req.TemplateId == "" || req.ActivityIntro == "" || req.PriceYuan == "" || req.SumNum < 0 {
		s.FailJsonExit(r, "缺少必要参数")
		return
	}
	if req.ActivityStartTime == nil || req.ActivityEndTime == nil || req.OpenAwardTime == nil || req.PayEndTime == nil {
		s.FailJsonExit(r, "缺少必要参数")
		return
	}
	if req.ActivityStartTime.Unix() > req.ActivityEndTime.Unix() {
		s.FailJsonExit(r, "开始时间不能大于结束时间")
		return
	}
	if req.ActivityEndTime.Unix() > req.OpenAwardTime.Unix() {
		s.FailJsonExit(r, "结束时间不能大于开奖时间")
		return
	}
	if req.OpenAwardTime.Unix() > req.PayEndTime.Unix() {
		s.FailJsonExit(r, "开奖时间不能大于支付截止时间")
		return
	}
	if req.ActivityType != 1 && req.ActivityType != 2 {
		s.FailJsonExit(r, "活动类型错误")
		return
	}
	if req.ActivityType == 2 {
		if req.GeneralBuyNum <= 0 {
			s.FailJsonExit(r, "抽签次数参数错误")
			return
		}
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
	req.Price = int(priceInt)
	req.Alias = utils.RandString(6)
	req.PublisherId = s.GetPublisherId(r)

	if req.CreatorId != 0 {
		creatorInfo, err := provider.Developer.GetCreatorInfo(req.CreatorId)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
		req.CreatorName = creatorInfo.Data.Name
		req.CreatorAvatar = creatorInfo.Data.LogoUrl
		req.CreatorNo = fmt.Sprintf("%06d", req.CreatorId)
	}
	if req.Id == 0 {
		err = service.AdminSubscribeActivity.Create(req.SubscribeActivity, req.Condition)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	} else {
		err = service.AdminSubscribeActivity.Update(req.SubscribeActivity, req.Condition)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	}
	s.SusJsonExit(r)
}

func (s *adminSubscribeActivity) List(r *ghttp.Request) {
	createStartTime := r.GetQueryString("createStartTime")
	createEndTime := r.GetQueryString("createEndTime")
	activityStartTimeA := r.GetQueryString("activityStartTimeA")
	activityStartTimeB := r.GetQueryString("activityStartTimeB")
	status := r.GetQueryString("status")
	activityEndTimeA := r.GetQueryString("activityEndTimeA")
	activityEndTimeB := r.GetQueryString("activityEndTimeB")
	activityType := r.GetQueryInt("activityType")
	searchVal := r.GetQueryString("searchVal")
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	publisherId := s.GetPublisherId(r)
	ret, err := service.AdminSubscribeActivity.ListByPage(activityType, publisherId, pageNum, createStartTime, createEndTime, activityStartTimeA, activityStartTimeB, status, activityEndTimeA, activityEndTimeB, searchVal, pageSize)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *adminSubscribeActivity) Detail(r *ghttp.Request) {
	id := r.GetQueryInt("id")
	if id == 0 {
		s.FailJsonExit(r, "参数错误")
		return
	}
	publisherId := s.GetPublisherId(r)
	ret, err := service.AdminSubscribeActivity.Detail(publisherId, id)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *adminSubscribeActivity) Delete(r *ghttp.Request) {
	id := r.GetInt("id")
	if id == 0 {
		s.FailJsonExit(r, "参数错误")
		return
	}
	count := service.SubscribeRecord.GetSubscribeRecordCount(id)
	if count > 0 {
		s.FailJsonExit(r, "活动已被认购，不能删除")
		return
	}
	publisherId := s.GetPublisherId(r)
	err := service.AdminSubscribeActivity.Delete(publisherId, id)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
}

func (s *adminSubscribeActivity) GetSubRecords(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 1)
	createdAtStart := r.GetQueryString("createdAtStart")
	createdAtEnd := r.GetQueryString("createdAtEnd")
	priceMinStr := r.GetQueryString("priceMinStr")
	priceMaxStr := r.GetQueryString("priceMaxStr")
	award := r.GetQueryInt("award")
	activityType := r.GetQueryInt("activityType")
	payStatus := r.GetQueryInt("payStatus")
	searchVal := r.GetQueryString("searchVal")
	priceMinStrValue, _ := decimal.NewFromString(priceMinStr)
	priceMinStrValue = priceMinStrValue.Mul(decimal.NewFromInt(100))
	priceMinInt := priceMinStrValue.IntPart()
	priceMaxStrValue, _ := decimal.NewFromString(priceMaxStr)
	priceMaxStrValue = priceMaxStrValue.Mul(decimal.NewFromInt(100))
	priceMaxInt := priceMaxStrValue.IntPart()
	publisherId := s.GetPublisherId(r)
	ret, err := service.AdminSubscribeActivity.GetSubRecords(activityType, pageNum, pageSize, publisherId, createdAtStart, createdAtEnd, int(priceMinInt), int(priceMaxInt), award, payStatus, searchVal)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *adminSubscribeActivity) Disable(r *ghttp.Request) {
	action := r.GetString("action")
	id := r.GetInt("id")
	publisherId := s.GetPublisherId(r)
	if action == "ON" {
		err := service.AdminSubscribeActivity.Disable(id, 0, publisherId)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	} else if action == "OFF" {
		err := service.AdminSubscribeActivity.Disable(id, 1, publisherId)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	}
	s.SusJsonExit(r)
}
