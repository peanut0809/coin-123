// ==========================================================================
// GFast自动生成model代码，无需手动修改，重新生成会自动覆盖.
// 生成日期：2022-05-13 11:20:13
// 生成路径: gfast/app/assets/model/blind_box_subscribe_records.go
// 生成人：DHX
// ==========================================================================

package model

import (
	"github.com/gogf/gf/os/gtime"
)

// BlindBoxSubscribeRecords is the golang structure for table blind_box_subscribe_records.
type SubscribeRecord struct {
	Id                 int         `orm:"id,primary" json:"id"` //
	ActivityType       int         `orm:"activity_type" json:"activityType"`
	Aid                int         `orm:"aid" json:"aid"`                 // 活动ID
	Name               string      `orm:"name" json:"name"`               // 盲盒名
	Icon               string      `orm:"icon" json:"icon"`               // 盲盒图
	UserId             string      `orm:"user_id" json:"userId"`          // 用户ID
	BuyNum             int         `orm:"buy_num" json:"buyNum"`          // 认购数量
	OrderNo            string      `orm:"order_no" json:"orderNo"`        // 订单号
	PayOrderNo         string      `orm:"pay_order_no" json:"payOrderNo"` //支付订单号
	SumPrice           int         `orm:"sum_price" json:"sumPrice"`      // 总价
	SumUnitCrystal     int         `orm:"sum_unit_crystal" json:"sumUnitCrystal"`
	SumUnitMonthTicket int         `orm:"sum_unit_month_ticket" json:"sumUnitMonthTicket"`
	SumUnitPrice       int         `orm:"sum_unit_price" json:"sumUnitPrice"`
	Award              int         `orm:"award" json:"award"` // 按中签类型
	AwardAt            *gtime.Time `orm:"award_at" json:"awardAt"`
	AwardNum           int         `orm:"award_num" json:"awardNum"` // 中签数量
	TicketType         string      `orm:"ticket_type" json:"ticketType"`
	PublishAsset       int         `orm:"publish_asset" json:"publishAsset"` // 发放状态
	PayStatus          int         `orm:"pay_status" json:"payStatus"`
	SubSum             int         `orm:"sub_sum" json:"subSum"`
	SubSumPeople       int         `orm:"sub_sum_people" json:"subSumPeople"`
	PublisherId        string      `orm:"publisher_id" json:"publisherId"`
	PayEndTime         *gtime.Time `orm:"pay_end_time" json:"payEndTime"`
	PaidAt             *gtime.Time `orm:"paid_at" json:"paidAt"`
	PayTicketMethod    string      `orm:"pay_ticket_method" json:"payTicketMethod"`
	PayMethod          string      `orm:"pay_method" json:"payMethod"`
	CreatedAt          *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt          *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}

type SubscribeRecordQueueData struct {
	SubscribeRecord
	FromUserId string `json:"fromUserId"`
	ToUserId   string `json:"toUserId"`
	TotalFee   int    `json:"totalFee"`
}

type SubscribeRecordList struct {
	List  []SubscribeRecord `json:"list"`
	Total int               `json:"total"`
}

type SubscribeRecordDetail struct {
	UserId       string      `json:"userId"`
	ActivityType int         `json:"activityType"`
	Name         string      `json:"name"`
	BuyNum       int         `json:"buyNum"`
	AwardNum     int         `json:"awardNum"`
	ConsumeUnit  string      `json:"consumeUnit"`
	Award        int         `json:"award"`
	AwardAt      *gtime.Time `json:"awardAt"`
	Icon         string      `json:"icon"`
	Aid          int         `json:"aid"`
}

type SubscribeListByOrderRetItem struct {
	BuyNum       int         `json:"buyNum"`
	SumPriceYuan string      `json:"sumPriceYuan"`
	SumPrice     int         `json:"sumPrice"`
	OrderNo      string      `json:"orderNo"` // 订单号
	Name         string      `json:"name"`
	Icon         string      `json:"icon"`
	Status       int         `json:"status"`
	PayOrderNo   string      `json:"payOrderNo"`
	PaidAt       *gtime.Time `json:"paidAt"`
	PayEndTime   *gtime.Time `json:"payEndTime"`
	PayMethod    string      `json:"payMethod"`
}

type SubscribeListByOrderRet struct {
	List  []SubscribeListByOrderRetItem `json:"list"`
	Total int                           `json:"total"`
}

type CreateOrderReq struct {
	OrderNo            string `json:"orderNo"` // 订单号
	SuccessRedirectUrl string `json:"successRedirectUrl"`
	ExitRedirectUrl    string `json:"exitRedirectUrl"`
}
