package api

import (
	"peanut-coin123/common"
	"peanut-coin123/model"
	"peanut-coin123/service"

	"github.com/gogf/gf/net/ghttp"
)

type front struct {
}

var Front = new(front)

func (s *front) CoinList(r *ghttp.Request) {
	var req model.CoinListReq
	err := r.Parse(&req)
	if err != nil {
		common.CommonMeans.ResponseFail(r, err.Error())
		return
	}
	cl, err2 := service.FrontService.CoinList(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r, err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, cl)
}

func (s *front) CoinIeoOffList(r *ghttp.Request) {
	var req model.CoinIeoOffReq
	err := r.Parse(&req)
	if err != nil {
		common.CommonMeans.ResponseFail(r, err.Error())
		return
	}
	cl, err2 := service.FrontService.CoinIeoOffList(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r, err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, cl)
}

func (s *front) CoinRenameList(r *ghttp.Request) {
	var req model.CoinRenameReq
	err := r.Parse(&req)
	if err != nil {
		common.CommonMeans.ResponseFail(r, err.Error())
		return
	}
	cl, err2 := service.FrontService.CoinRenameList(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r, err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, cl)
}

func (s *front) CoinCoreList(r *ghttp.Request) {
	var req model.CoinCoreReq
	err := r.Parse(&req)
	if err != nil {
		common.CommonMeans.ResponseFail(r, err.Error())
		return
	}
	cl, err2 := service.FrontService.CoinCoreList(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r, err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, cl)
}
