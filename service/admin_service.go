package service

import (
	"fmt"
	"peanut-coin123/model"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
)

type adminService struct{}

var AdminService = new(adminService)

// 创建代币上架列表
func (s *adminService) Create(in model.CoinItems) (err error) {
	// 获取详情
	var coinItems *model.CoinItems
	m := g.DB().Model("coin123_items")
	err = m.Where("coin_name", in.CoinName).Scan(&coinItems)
	if err != nil {
		return
	}
	if coinItems != nil {
		err = fmt.Errorf(in.CoinName + "数据已经存在")
		return
	}
	r, err2 := g.DB().Model("coin123_items").Insert(&in)
	if err2 != nil {
		return
	}
	r.LastInsertId()
	return
}

// 创建代币上架列表
func (s *adminService) CoinList(in model.CoinListReq) (list model.CoinList, err error) {
	m := g.DB().Model("coin123_items")
	if in.CoinName != "" {
		m = m.Where("coin_name", in.CoinName)
	}
	if in.Type > 0 {
		m = m.Where("type", in.Type)
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

// 创建核心公告列表
func (s *adminService) CoreContextCreat(in model.CoinCoreContext) (err error) {
	r, err2 := g.DB().Model("coin123_core_context").Insert(&in)
	if err2 != nil {
		return
	}
	r.LastInsertId()
	return
}

// 创建代币上架列表
func (s *adminService) CoinCoreList(in model.CoinCoreReq) (list model.CoinCoreList, err error) {
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

// 创建ieo/下架列表
func (s *adminService) IeoOffCreat(in model.CoinIeoOffItems) (err error) {
	r, err2 := g.DB().Model("coin123_ieo_off_items").Insert(&in)
	if err2 != nil {
		return
	}
	r.LastInsertId()
	return
}
func (s *adminService) CoinIeoOffList(in model.CoinIeoOffReq) (list model.CoinIeoOffList, err error) {
	m := g.DB().Model("coin123_ieo_off_items")
	if in.CoinName != "" {
		m = m.Where("coin_name", in.CoinName)
	}
	if in.Type != "" {
		m = m.Where("type", in.Type)
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

// 创建改名代币
func (s *adminService) CoinRenameCreat(in model.CoinRenameItems) (err error) {
	r, err2 := g.DB().Model("coin123_rename_items").Insert(&in)
	if err2 != nil {
		return
	}
	r.LastInsertId()
	return
}
func (s *adminService) CoinRenameList(in model.CoinRenameReq) (list model.CoinRenameList, err error) {
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
