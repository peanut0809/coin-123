package service

import (
	"peanut-coin123/model"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
)

type frontService struct{}

var FrontService = new(frontService)

func (s *frontService) CoinList(in model.CoinListReq) (list model.CoinList, err error) {
	m := g.DB().Model("coin123_items")
	if in.CoinName != "" {
		m = m.Where("coin_name", in.CoinName)
	}
	if in.Type > 0 {
		m = m.Where("type", in.Type)
	}
	if in.DateContext != "" {
		m = m.Where("date_context", in.DateContext)
	}
	total, err := m.Count()
	if err != nil {
		err = gerror.New("获取总行数失败")
		return
	}
	list.Total = total
	items := make([]*model.CoinItems, 0)
	err = m.Order("id DESC").Page(in.Page, in.PageSize).Scan(&items)
	if err != nil {
		return
	}
	list.List = items
	if err != nil {
		return
	}
	return
}
func (s *frontService) CoinIeoOffList(in model.CoinIeoOffReq) (list model.CoinIeoOffList, err error) {
	m := g.DB().Model("coin123_ieo_off_items")
	if in.CoinName != "" {
		m = m.Where("coin_name", in.CoinName)
	}
	if in.Type != "" {
		m = m.Where("type", in.Type)
	}
	if in.DateContext != "" {
		m = m.Where("date_context", in.DateContext)
	}
	total, err := m.Count()
	if err != nil {
		err = gerror.New("获取总行数失败")
		return
	}
	list.Total = total
	items := make([]*model.CoinIeoOffItems, 0)
	err = m.Order("id DESC").Page(in.Page, in.PageSize).Scan(&items)
	if err != nil {
		return
	}
	list.List = items
	if err != nil {
		return
	}
	return
}

func (s *frontService) CoinRenameList(in model.CoinRenameReq) (list model.CoinRenameList, err error) {
	m := g.DB().Model("coin123_rename_items")
	total, err := m.Count()
	if err != nil {
		err = gerror.New("获取总行数失败")
		return
	}
	list.Total = total
	items := make([]*model.CoinRenameItems, 0)
	err = m.Order("id DESC").Page(in.Page, in.PageSize).Scan(&items)
	if err != nil {
		return
	}
	list.List = items
	if err != nil {
		return
	}
	return
}

func (s *frontService) CoinCoreList(in model.CoinCoreReq) (list model.CoinCoreList, err error) {
	m := g.DB().Model("coin123_core_context")
	total, err := m.Count()
	if err != nil {
		err = gerror.New("获取总行数失败")
		return
	}
	list.Total = total
	items := make([]*model.CoinCoreContext, 0)
	err = m.Order("id DESC").Page(in.Page, in.PageSize).Scan(&items)
	if err != nil {
		return
	}
	list.List = items
	if err != nil {
		return
	}
	return
}
