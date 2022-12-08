package api

import (
	"meta_launchpad/model"
	"meta_launchpad/service"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"github.com/shopspring/decimal"
)

type adminEquity struct {
	api.CommonBase
}

var AdminEquity = new(adminEquity)

func (s *adminEquity) Create(r *ghttp.Request) {
	var req model.CreateEquityActivityReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}

	if req.Name == "" || req.TimeType <= 0 || req.SubLimitType <= 0 || req.TemplateId == "" || req.Number <= 0 || req.LimitBuy <= 0 || req.LimitType <= 0 || req.ActivityStartTime == nil || req.ActivityEndTime == nil {
		s.FailJsonExit(r, "参数错误")
		return
	}
	if req.LimitType == model.EQUITY_ACTIVITY_LIMIT_TYPE1 {
		if req.LimitType == 1 {
			req.LimitBuy = 999
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
	if req.ActivityStartTime.Unix() > req.ActivityEndTime.Unix() {
		s.FailJsonExit(r, "活动开始时间不能大于结束时间")
		return
	}
	req.Price = int(priceInt)

	err = service.AdminEquity.Create(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
}

// 用户导入解析
func (s *adminEquity) Import(r *ghttp.Request) {
	var req model.CreateEquityActivityReq
	req.IsCreate = false
	result, err := service.AdminEquity.HandelExcelUser(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, result)
}

// 获取详情
func (s *adminEquity) Item(r *ghttp.Request) {
	templateId := r.GetQueryString("templateId")
	if templateId == "" {
		s.FailJsonExit(r, "活动标识为空")
		return
	}
	ret, err := service.AdminEquity.Item(templateId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

//Invalid 失效下架活动
// 获取详情
func (s *adminEquity) Invalid(r *ghttp.Request) {
	equityId := r.GetQueryInt("equityId")
	if equityId <= 0 {
		s.FailJsonExit(r, "活动标识为空")
		return
	}
	err := service.AdminEquity.Invalid(equityId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
}

// 获取专属活动用户明细
func (s *adminEquity) UserItems(r *ghttp.Request) {

	var req model.EquityUserReq
	if err := r.Parse(&req); err != nil {
		s.FailJsonExit(r, err.(gvalid.Error).FirstString())
		return
	}
	//获取参数
	req.PublisherId = s.GetPublisherId(r)
	if req.EquityId <= 0 {
		s.FailJsonExit(r, "权益活动id异常")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	ret, err := service.AdminEquity.UserItems(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

// 获取发行商发售权益活动明细
func (s *adminEquity) EquityItems(r *ghttp.Request) {
	var req model.AdminEquityReq
	if err := r.Parse(&req); err != nil {
		s.FailJsonExit(r, err.(gvalid.Error).FirstString())
		return
	}
	//获取参数
	req.PublisherId = s.GetPublisherId(r)

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	ret, err := service.AdminEquity.EquityActivityItems(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

// 订单列表
func (s *adminEquity) OrderItems(r *ghttp.Request) {
	var req model.AdminEquityOrderReq
	if err := r.Parse(&req); err != nil {
		s.FailJsonExit(r, err.(gvalid.Error).FirstString())
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	req.PublisherId = s.GetPublisherId(r)
	ret, err := service.AdminEquity.OrderItems(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

// 订单导出
func (s *adminEquity) OrderExport(r *ghttp.Request) {

}
