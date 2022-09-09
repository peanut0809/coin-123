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
	publisherId := c.GetPublisherId(r)
	var params model.BannerReq
	if err := r.Parse(&params); err != nil {
		c.FailJsonCodeExit(r, err)
		return
	}
	total, page, list, err := service.Banner.List(publisherId, params)
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

// Create 新增，修改
func (c *banner) Create(r *ghttp.Request) {
	var params model.BannerCreateReq
	if err := r.Parse(&params); err != nil {
		c.FailJsonCodeExit(r, err)
		return
	}
	params.PublisherId = c.GetPublisherId(r)
	state, err := service.Banner.Create(params)
	if err != nil {
		c.FailJsonCodeExit(r, err)
	}
	if state != "" {
		c.SusJsonExit(r, state)
	}
	if params.Id == 0 {
		c.SusJsonExit(r, "添加成功")
	} else {
		c.SusJsonExit(r, "修改成功")
	}
}

// Delete 修改
func (c *banner) Delete(r *ghttp.Request) {
	id := r.GetInt("id")

	state, err := service.Banner.Delete(id)
	if err != nil {
		c.FailJsonCodeExit(r, err)
	}
	if state != "" {
		c.SusJsonExit(r, state)
	}
	c.SusJsonExit(r, "删除成功")
}

// StateEdit 修改
func (c *banner) StateEdit(r *ghttp.Request) {

	id := r.GetInt("id")
	state := r.GetInt("state")
	err := service.Banner.StateEdit(id, state)
	if err != nil {
		c.FailJsonCodeExit(r, err)
	}
	c.SusJsonExit(r, "状态修改成功")
}

// GetFrontList 前段展示
func (c *banner) GetFrontList(r *ghttp.Request) {
	//publisherId := c.GetPublisherId(r)
	publisherId := r.GetString("publisherId")
	if publisherId == "" {
		c.FailJsonExit(r, "发行商ID不能为空")
	}
	list := service.Banner.FrontList(publisherId)
	req := g.Map{
		"list": list,
	}
	c.SusJsonExit(r, req)
}

// GetRichText 前段展示
func (c *banner) GetRichText(r *ghttp.Request) {
	//publisherId := c.GetPublisherId(r)
	id := r.GetInt("id")
	if id == 0 {
		c.FailJsonExit(r, "id不能为空")
	}
	list := service.Banner.RichText(id)
	c.SusJsonExit(r, list.JumpUrl)
}
