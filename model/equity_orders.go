package model

import (
	"github.com/gogf/gf/os/gtime"
)

const WAIT_PAY = 1 // 待支付
const PAID = 2     // 已支付
const TIMEOUT = 3  // 已超时
const CANCEL = 4   // 已取消

type EquityOrder struct {
	Id           int         `orm:"id" json:"id"`                      // 主键
	PublisherId  string      `orm:"publisher_id" json:"publisherId"`   // 发行商ID
	OrderNo      string      `orm:"order_no" json:"orderNo"`           // 订单号
	Num          int         `orm:"num" json:"num"`                    // 购买数量
	RealFee      int         `orm:"real_fee" json:"realFee"`           // 实际成交额
	ActivityId   int         `orm:"activity_id" json:"activityId"`     // 活动ID
	UserId       string      `orm:"user_id" json:"userId"`             // 用户ID
	ActivityName string      `orm:"activity_name" json:"activityName"` // 活动名字
	UserName     string      `orm:"user_name" json:"userName"`         // 用户名字
	UserPhone    string      `orm:"user_phone" json:"userPhone"`       // 用户手机号
	Icon         string      `orm:"icon" json:"icon"`                  // icon
	Status       int         `orm:"status" json:"status"`              // 1.待支付；2.已支付；3.已超时；4.已取消
	Price        int         `orm:"price" json:"price"`                // 单价
	PayAt        *gtime.Time `orm:"pay_at" json:"payAt"`               // 支付时间
	PayMethod    string      `orm:"pay_method" json:"payMethod"`       // 支付方式
	PayExpireAt  *gtime.Time `orm:"pay_expire_at" json:"payExpireAt"`  // 支付过期时间
	CreatedAt    *gtime.Time `orm:"created_at" json:"createdAt"`       // 新建时间
	UpdatedAt    *gtime.Time `orm:"updated_at" json:"updatedAt"`       // 更新时间
	LimitType    int         `orm:"limit_type" json:"limitType"`       // 限购类型 1 按每人限购 2 白名单限购
}

type EquityOrderFull struct {
	*EquityOrder
	PriceYuan   string `json:"priceYuan"`
	RealFeeYuan string `json:"realFeeYuan"`
	LastSec     int64  `json:"lastSec"`
}

type EquityOrderList struct {
	List  []*EquityOrderFull
	Total int
}

type AdminEquityOrderReq struct {
	PublisherId string `json:"publisherId"`
	ActivityId  int    `json:"activityId"`
	Page        int    `json:"pageNum"`
	PageSize    int    `json:"pageSize"`
	Phone       int    `json:"phone"`
	Status      int    `json:"status"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	MinPrice    int    `json:"minPrice"`
	MaxPrice    int    `json:"maxPrice"`
	OrderNo     string `json:"orderNo"`
}
