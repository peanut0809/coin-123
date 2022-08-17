package model

import (
	"github.com/gogf/gf/os/gtime"
)

type Banner struct {
	Id           int         `json:"id" orm:"id"`
	Name         string      `json:"name" orm:"name"`                   // banner 名称
	Remarks      string      `json:"remarks" orm:"remarks"`             // 备注
	Image        string      `json:"image" orm:"image"`                 // 图片地址
	JumpType     string      `json:"jumpType" orm:"jump_type"`          // 跳转链接类型
	JumpUrl      string      `json:"jumpUrl" orm:"jump_url"`            // 跳转链接地址
	Sort         int         `json:"sort" orm:"sort"`                   // 排序
	State        int         `json:"state" orm:"state"`                 // 状态 0：未上架 1：已上架 2：已下架
	TimingState  int         `json:"timingState" orm:"timing_state"`    // 定时任务
	GoodsOnTime  *gtime.Time `json:"goodsOnTime" orm:"goods_on_time"`   // 上架时间
	GoodsOffTime *gtime.Time `json:"goodsOffTime" orm:"goods_off_time"` // 下架时间
	CreatedAt    string      `json:"createdAt" orm:"created_at"`        // 创建时间
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
	PublisherId  string `json:"publisher_id"`
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
