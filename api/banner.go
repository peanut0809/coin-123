package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

type banner struct {
	api.CommonBase
}

var Banner = new(banner)

func (c *banner) GetList(r *ghttp.Request) {
	var params model.BannerReq
	if err := r.Parse(&params); err != nil {
		c.FailJsonCodeExit(r, err)
		return
	}
	total, page, list, err := service.Banner.List(params)
	if err != nil {
		c.FailJsonExit(r, "获取数据失败")
	}
	req := g.Map{
		"total": total,
		"list":  list,
		"page":  page,
	}
	c.SusJsonExit(r, req)
}

// Add 新增
func (c *banner) Add(r *ghttp.Request) {
	var params model.BannerAddReq
	if err := r.Parse(&params); err != nil {
		c.FailJsonCodeExit(r, err)
		return
	}
	err := service.Banner.Add(params)
	if err != nil {
		c.FailJsonCodeExit(r, err)
	}
	c.SusJsonExit(r, "添加成功")
}

// Edit 修改
func (c *banner) Edit(r *ghttp.Request) {
	var params model.BannerEditReq
	if err := r.Parse(&params); err != nil {
		c.FailJsonCodeExit(r, err)
		return
	}
	err := service.Banner.Edit(params)
	if err != nil {
		c.FailJsonCodeExit(r, err)
	}
	c.SusJsonExit(r, "修改成功")
}

// Delete 修改
func (c *banner) Delete(r *ghttp.Request) {
	id := r.GetInt("id")
	err := service.Banner.Delete(id)
	if err != nil {
		c.FailJsonCodeExit(r, err)
	}
	c.SusJsonExit(r, "删除成功")
}
