package model

import (
	"github.com/gogf/gf/os/gtime"
)

type SeckillOrder struct {
	Id          int         `orm:"id" json:"id"` // 主键
	PublisherId string      `orm:"publisher_id" json:"publisherId"`
	OrderNo     string      `orm:"order_no" json:"orderNo"`          // 订单号
	Num         int         `orm:"num" json:"num"`                   // 购买数量
	RealFee     int         `orm:"real_fee" json:"realFee"`          // 实际成交额
	UserId      string      `orm:"user_id" json:"userId"`            // 用户ID
	Aid         int         `orm:"aid" json:"aid"`                   // 活动ID
	Name        string      `orm:"name" json:"name"`                 // 名字
	Icon        string      `orm:"icon" json:"icon"`                 // icon
	Status      int         `orm:"status" json:"status"`             // 1.待支付；2.已支付；3.已超时；4.已取消
	Price       int         `orm:"price" json:"price"`               // 单价
	PayAt       *gtime.Time `orm:"pay_at" json:"payAt"`              // 支付时间
	PayMethod   string      `orm:"pay_method" json:"payMethod"`      // 支付方式
	PayExpireAt *gtime.Time `orm:"pay_expire_at" json:"payExpireAt"` // 支付过期时间
	CreatedAt   *gtime.Time `orm:"created_at" json:"createdAt"`      // 新建时间
	UpdatedAt   *gtime.Time `orm:"updated_at" json:"updatedAt"`      // 更新时间
}

type SeckillOrderFull struct {
	SeckillOrder
	PriceYuan   string `json:"priceYuan"`
	RealFeeYuan string `json:"realFeeYuan"`
	LastSec     int64  `json:"lastSec"`
}

type SeckillOrderList struct {
	Total int                `json:"total"`
	List  []SeckillOrderFull `json:"list"`
}

type AdminSeckillOrderFull struct {
	SeckillOrder
	StatusTxt   string `json:"statusTxt"`
	UserName    string `json:"userName"`
	UserPhone   string `json:"userPhone"`
	RealFeeYuan string `json:"realFeeYuan"`
}

type AdminSeckillOrderByPage struct {
	Total int                     `json:"total"`
	List  []AdminSeckillOrderFull `json:"list"`
}
