package api

import (
	"meta_launchpad/model"
	"meta_launchpad/service"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
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

	if req.Name == "" || req.Number <= 0 || req.LimitBuy <= 0 || req.LimitType <= 0 || req.ActivityStartTime == nil || req.ActivityEndTime == nil {
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
	result, err := service.AdminEquity.HandelExcelUser(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, result)
}
