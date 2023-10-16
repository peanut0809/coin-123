package admin

import (
	"peanut-coin123/common"
	"peanut-coin123/model"
	"peanut-coin123/service"

	"github.com/gogf/gf/net/ghttp"
)

// 后台管理
type admin struct {
}

var Admin = new(admin)

// 创建币安上币交易对
func (s *admin) CreateCoinItems(r *ghttp.Request) {
	var req model.CoinItems
	err := r.Parse(&req)
	if err != nil {
		common.CommonMeans.ResponseFail(r, err.Error())
		return
	}
	if req.CoinName == "" {
		common.CommonMeans.ResponseFail(r, "coin name is empty")
		return
	}
	if req.Context == "" {
		common.CommonMeans.ResponseFail(r, "coin context is empty")
		return
	}
	if req.OriginCreatedAt == "" {
		common.CommonMeans.ResponseFail(r, "coin origin create at is empty")
		return
	}
	if req.DateContext == "" {
		common.CommonMeans.ResponseFail(r, "coin date context create at is empty")
		return
	}
	if req.Type <= 0 {
		common.CommonMeans.ResponseFail(r, "spots contract err")
		return
	}
	err2 := service.AdminService.Create(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r, err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, nil)
}
func (s *admin) CoinList(r *ghttp.Request) {
	common.CommonMeans.ResponseSuccess(r, nil)
}

// 创建币安IEO/下架交易对
func (s *admin) CreateIeoOffCoin(r *ghttp.Request) {
	var req model.CoinIeoOffItems
	err := r.Parse(&req)
	if err != nil {
		common.CommonMeans.ResponseFail(r, err.Error())
		return
	}
	if req.Context == "" {
		common.CommonMeans.ResponseFail(r, "请输入内容")
		return
	}
	if req.OriginCreatedAt == "" {
		common.CommonMeans.ResponseFail(r, "请输入OriginCreatedAt")
		return
	}
	if req.Type == "" {
		common.CommonMeans.ResponseFail(r, "请选择ieo/off")
		return
	}
	err2 := service.AdminService.IeoOffCreat(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r, err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, nil)
}

// 创建币安核心公告
func (s *admin) CreateCoreContext(r *ghttp.Request) {
	var req model.CoinCoreContext
	err := r.Parse(&req)
	if err != nil {
		common.CommonMeans.ResponseFail(r, err.Error())
		return
	}
	if req.Context == "" || req.OriginCreatedAt == "" {
		common.CommonMeans.ResponseFail(r, "内容异常")
		return
	}
	err2 := service.AdminService.CoreContextCreat(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r, err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, nil)
}

// 创建币安改名
func (s *admin) CreateCoinRename(r *ghttp.Request) {
	var req model.CoinRenameItems
	err := r.Parse(&req)
	if err != nil {
		common.CommonMeans.ResponseFail(r, err.Error())
		return
	}
	if req.AfterName == "" || req.OriginalName == "" {
		common.CommonMeans.ResponseFail(r, "名字为空")
		return
	}
	if req.Context == "" || req.OriginCreatedAt == "" {
		common.CommonMeans.ResponseFail(r, "请输入内容/文章时间")
		return
	}
	err2 := service.AdminService.CoinRenameCreat(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r, err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, nil)
}

// 列表集合

func (s *admin) CoinItems(r *ghttp.Request) {

}
func (s *admin) CoinCoreItems(r *ghttp.Request) {

}
func (s *admin) CoinIeoOffCoreItems(r *ghttp.Request) {

}
func (s *admin) CoinRenameCoreItems(r *ghttp.Request) {

}
