package service

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"time"
)

type banner struct {
}

var Banner = new(banner)

// List 数据获取，条件筛选
func (c *banner) List(params model.BannerReq) (total, page int, list []model.Banner, err error) {
	db := g.DB().Model("banner")
	if params.CreatedStart != "" && params.CreatedEnd != "" {
		db = db.Where("created_at >= ? and created_at <= ?", params.CreatedStart, params.CreatedEnd)
	}
	if params.GoodsOnStart != "" && params.GoodsOnEnd != "" {
		db = db.Where("goods_on_time >= ? and goods_on_time <= ?", params.GoodsOnStart, params.GoodsOnEnd)
	}
	if params.GoodsOffStart != "" && params.GoodsOffEnd != "" {
		db = db.Where("goods_off_time >= ? and goods_off_time <= ?", params.GoodsOffStart, params.GoodsOffEnd)
	}
	if params.State != "" {
		db = db.Where("state = ?", params.State)
	}
	if params.Name != "" {
		db = db.Where("`name` like ?", "%"+params.Name+"%")
	}
	total, err = db.Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}
	if params.PageNum == 0 {
		params.PageNum = 1
	}
	page = params.PageNum
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	err = db.Page(page, params.PageSize).Scan(&list)
	if err != nil {
		g.Log().Error(err)
		return
	}
	return
}

// Add 新增
func (c *banner) Add(params model.BannerAddReq) (err error) {
	params.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	_, err = g.DB().Model("banner").Insert(params)
	return
}

// Edit 修改
func (c *banner) Edit(params model.BannerEditReq) (err error) {
	params.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	_, err = g.DB().Model("banner").Where("id = ?", params.Id).Update(params)
	return
}

// Delete 修改
func (c *banner) Delete(id int) (state string, err error) {
	var list model.Banner
	err = g.DB().Model("banner").Where("id = ?", id).Scan(&list)
	if err != nil {
		return
	}

	if list.State == 1 {
		state = "正在上架中无法删除"
	} else {
		_, err = g.DB().Model("banner").Delete("id = ?", id)
	}
	return
}
