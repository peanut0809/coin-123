package model

import (
	"github.com/gogf/gf/os/gtime"
)

const EQUITY_ACTIVITY_STATUS1 = 1 // 上架
const EQUITY_ACTIVITY_STATUS2 = 2 // 下架

const EQUITY_ACTIVITY_LIMIT_TYPE1 = 1 // 每人限购
const EQUITY_ACTIVITY_LIMIT_TYPE2 = 2 // 专属限购 白名单用户

const SubSetEquityResultKey = "meta_launchpad:activity_equity_result:%s"
const EquityActivityStatusWait = 0
const EquityActivityStatusIng = 1
const EquityActivityStatusEnd = 2

type EquityActivity struct {
	Id                int         `orm:"id,primary" json:"id"`                         // 活动iD
	PublisherId       string      `orm:"publisher_id" json:"publisherId"`              // 发行商ID
	Name              string      `orm:"name" json:"name"`                             // 活动名称
	Price             int         `orm:"price" json:"price"`                           // 发售价,单位：分
	ActivityStartTime *gtime.Time `orm:"activity_start_time" json:"activityStartTime"` // 活动开始时间
	ActivityEndTime   *gtime.Time `orm:"activity_end_time" json:"activityEndTime"`     // 活动结束时间
	LimitBuy          int         `orm:"limit_buy" json:"limitBuy"`                    // 限购类型 1 按每人限购 2 白名单限购  1 每人限购数量
	LimitType         int         `orm:"limit_type" json:"limitType"`                  // 限购类型 1 按每人限购 2 白名单限购
	Number            int         `orm:"number" json:"number"`                         // 总数量
	Status            int         `orm:"status" json:"status"`                         // 活动状态1:上架 2:下架
	CreatedAt         *gtime.Time `orm:"created_at" json:"createdAt"`                  // 新建时间
	UpdatedAt         *gtime.Time `orm:"updated_at" json:"updatedAt"`                  // 更新时间
}

type CreateWhiteActivityReq struct {
	EquityActivity
	PriceYuan string `json:"priceYuan"`
}

type EquityActivityList struct {
	List  []*EquityActivity `json:"list"`
	Total int               `json:"total"`
}

type EquitySubResult struct {
	Reason  string `json:"reason"`
	OrderNo string `json:"orderNo"`
	Type    string `json:"type"`
	Step    string `json:"step"`
}

type EquityOrderReq struct {
	Id                 int    `json:"id"`
	Num                int    `json:"num"`
	ClientIp           string `json:"clientIp"`
	SuccessRedirectUrl string `json:"successRedirectUrl"`
	ExitRedirectUrl    string `json:"exitRedirectUrl"`
	PublisherId        string `json:"publisherId"`
	PlatformAppId      string `json:"platformAppId"`
	UserId             string `json:"userId"`
	OrderNo            string `json:"orderNo"`
}
