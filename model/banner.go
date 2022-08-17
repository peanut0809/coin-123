package model

import (
	"github.com/gogf/gf/os/gtime"
)

type Banner struct {
	Id           int         `json:"id" orm:"id"`
	Name         string      `json:"name" orm:"name"`
	Remarks      string      `json:"remarks" orm:"remarks"`
	Image        string      `json:"image" orm:"image"`
	JumpType     string      `json:"jumpType" orm:"jump_type"`
	JumpUrl      string      `json:"jumpUrl" orm:"jump_url"`
	Sort         int         `json:"sort" orm:"sort"`
	State        int         `json:"state" orm:"state"`
	TimingState  int         `json:"timingState" orm:"timing_state"`
	GoodsOnTime  *gtime.Time `json:"goodsOnTime" orm:"goods_on_time"`
	GoodsOffTime *gtime.Time `json:"goodsOffTime" orm:"goods_off_time"`
	CreatedAt    string      `json:"createdAt" orm:"created_at"`
}

// BannerReq 请求字段
type BannerReq struct {
	PageNum       int    `json:"pageNum"`
	PageSize      int    `json:"pageSize"`
	CreatedStart  string `json:"createdStart"`
	CreatedEnd    string `json:"createdEnd"`
	GoodsOnStart  string `json:"goodsOnStart"`
	GoodsOnEnd    string `json:"goodsOnEnd"`
	GoodsOffStart string `json:"goodsOffStart"`
	GoodsOffEnd   string `json:"goodsOffEnd"`
	State         string `json:"state"`
	Name          string `json:"name"`
}

// BannerCreateReq 新增、修改
type BannerCreateReq struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Remarks      string `json:"remarks"`
	Image        string `json:"image"`
	JumpType     string `json:"jumpType"`
	JumpUrl      string `json:"jumpUrl"`
	Sort         int    `json:"sort"`
	TimingState  int    `json:"timingState"`
	GoodsOnTime  string `json:"goodsOnTime"`
	GoodsOffTime string `json:"goodsOffTime"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
