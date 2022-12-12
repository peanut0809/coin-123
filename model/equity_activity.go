package model

import (
	"meta_launchpad/provider"

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
const EQUITY_LIMITBUY = 999

type EquityActivity struct {
	Id                int         `orm:"id,primary" json:"id"`                         // 活动iD
	PublisherId       string      `orm:"publisher_id" json:"publisherId"`              // 发行商ID
	TemplateId        string      `orm:"template_id" json:"templateId"`                // 模板ID
	AppId             string      `orm:"app_id" json:"appId"`                          // 应用ID
	Name              string      `orm:"name" json:"name"`                             // 活动名称
	Price             int         `orm:"price" json:"price"`                           // 发售价,单位：分
	ActivityStartTime *gtime.Time `orm:"activity_start_time" json:"activityStartTime"` // 活动开始时间
	ActivityEndTime   *gtime.Time `orm:"activity_end_time" json:"activityEndTime"`     // 活动结束时间
	ActivityStatus    int         `json:"activityStatus"`                              // 活动状态
	ActivityStatusTxt string      `json:"activityStatusTxt"`                           // 活动状态中文
	LimitBuy          int         `orm:"limit_buy" json:"limitBuy"`                    // 限购类型 1 按每人限购 2 白名单限购  1 每人限购数量
	LimitType         int         `orm:"limit_type" json:"limitType"`                  // 限购类型 1 按每人限购 2 白名单限购
	TimeType          int         `orm:"time_type" json:"timeType"`                    // 时间类型 1 立即上架 2 按自定义时间
	SubLimitType      int         `orm:"sub_limit_type" json:"subLimitType"`           // 子类型限购 当limit_type = 1时候 分 不限购最大999/限购数量
	Number            int         `orm:"number" json:"number"`                         // 库存
	TotalNumber       int         `orm:"total_number" json:"totalNumber"`              // 总库存
	CoverImgUrl       string      `orm:"cover_img_url" json:"coverImgUrl"`
	Status            int         `orm:"status" json:"status"`        // 活动状态1:上架 2:下架
	CreatedAt         *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt         *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
	NfrSec            int         `orm:"nfr_sec" json:"nfrSec"`
}

type EquityActivityFull struct {
	EquityActivity
	LastSec         int64                       `json:"lastSec"`
	PriceYuan       string                      `json:"priceYuan"`
	AssetCateString string                      `json:"assetCateString"`
	AssetTotal      int                         `json:"assetTotal"`
	AssetCreateAt   *gtime.Time                 `json:"assetCreateAt"`
	AssetDetailImg  string                      `json:"assetDetailImg"`
	AssetPic        string                      `json:"assetPic"`
	NfrDay          int                         `json:"nfrDay"`
	ChainName       string                      `json:"chainName"`
	ChainAddr       string                      `json:"chainAddr"`
	ChainType       int                         `json:"chainType"`
	CopyrightInfo   []provider.TplCopyrightInfo `json:"copyrightInfo"`
}

type CreateEquityActivityReq struct {
	EquityActivity
	PriceYuan string `json:"priceYuan"`
	ExcelFile string `json:"excelFile"` //导入名单集合
	IsCreate  bool   `json:"isCreate"`
}

type AdminEquityReq struct {
	PublisherId string `json:"publisherId"`
	TemplateId  string `json:"templateId"`
	Page        int    `json:"pageNum"`
	PageSize    int    `json:"pageSize"`
	Name        string `json:"name"`
	Status      int    `json:"status"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
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

type ImportItems struct {
	ErrItems    []ImportItem
	SuccItems   []ImportItem
	HaveErr     bool
	Total       int
	Number      int
	AssetsCount int
}
type ImportItem struct {
	Phone      string `orm:"phone" json:"phone"`
	LimitNum   int    `orm:"limit_num" json:"limitNum"`
	ErrMessage string
	UserId     string `orm:"user_id" json:"userId"`
}

type AssetItems struct {
	Count int `json:"count"` //
}

type EquityPutItems struct {
	TemplateId string `orm:"templateId" json:"templateId"`
	Status     int    `orm:"status" json:"status"`
}
