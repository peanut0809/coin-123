package model

import (
	"meta_launchpad/provider"

	"github.com/gogf/gf/os/gtime"
)

type CreateSubscribeActivityReq struct {
	SubscribeActivity
	PriceYuan string               `json:"priceYuan"`
	Condition []SubscribeCondition `json:"condition"`
}

type SubscribeActivity struct {
	Id                int         `orm:"id,primary" json:"id"` // 活动iD
	ActivityType      int         `orm:"activity_type" json:"activityType"`
	Alias             string      `orm:"alias" json:"alias"`                           // 别名
	Name              string      `orm:"name" json:"name"`                             // 活动名
	ActivityStartTime *gtime.Time `orm:"activity_start_time" json:"activityStartTime"` // 活动开始时间
	ActivityEndTime   *gtime.Time `orm:"activity_end_time" json:"activityEndTime"`     //活动结束时间
	StartTime         *gtime.Time `orm:"start_time" json:"startTime"`                  // 普通用户认购开始时间
	Price             int         `orm:"price" json:"price"`                           // 发售价,单位：分
	TicketInfo        string      `orm:"ticket_info" json:"ticketInfo"`
	AssetIntro        string      `orm:"asset_intro" json:"assetIntro"`
	ActivityIntro     string      `orm:"activity_intro" json:"activityIntro"`
	CoverImgUrl       string      `orm:"cover_img_url" json:"coverImgUrl"` //商品封面图
	PublisherId       string      `orm:"publisher_id" json:"publisherId"`
	AppId             string      `orm:"app_id" json:"appId"`           // 应用ID
	AssetType         string      `orm:"asset_type" json:"assetType"`   // 资产类型
	AssetType2        string      `orm:"asset_type2" json:"assetType2"` // 资产类型
	AssetType3        string      `orm:"asset_type3" json:"assetType3"` // 资产类型
	TemplateId        string      `orm:"template_id" json:"templateId"` // 模板ID
	SumNum            int         `orm:"sum_num" json:"sumNum"`         // 总数
	RemainNum         int         `orm:"remain_num" json:"remainNum"`
	OpenAwardTime     *gtime.Time `orm:"open_award_time" json:"openAwardTime"` // 普通中签公布时间
	PayEndTime        *gtime.Time `orm:"pay_end_time" json:"payEndTime"`       //发放资产时间
	AwardStatus       int         `orm:"award_status" json:"awardStatus"`      // 开奖状态
	SubSum            int         `orm:"sub_sum" json:"subSum"`
	SubSumPeople      int         `orm:"sub_sum_people" json:"subSumPeople"`
	GeneralBuyNum     int         `orm:"general_buy_num" json:"generalBuyNum"`
	GeneralNumMethod  int         `orm:"general_num_method" json:"generalNumMethod"`
	AwardMethod       int         `orm:"award_method" json:"awardMethod"`
	AwardCompleteTime *gtime.Time `orm:"award_complete_time" json:"awardCompleteTime"`
	Disable           int         `orm:"disable" json:"disable"`
	NfrSec            int         `orm:"nfr_sec" json:"nfrSec"`
	CreatorId         int         `orm:"creator_id" json:"creatorId"`
	CreatorName       string      `orm:"creator_name" json:"creatorName"`
	CreatorAvatar     string      `orm:"creator_avatar" json:"creatorAvatar"`
	CreatorNo         string      `orm:"creator_no" json:"creatorNo"`
	CreatedAt         *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt         *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}

type TicketInfoJson struct {
	Use       bool   `json:"use"`
	Type      string `json:"type"`
	Num       int    `json:"num"`
	MaxBuyNum int    `json:"maxBuyNum"`
	UnitNum   int    `json:"unitNum"`
	IsShare   int    `json:"isShare"`
}

type SubscribeActivityFull struct {
	SumNum        int                         `json:"sumNum"`
	CoverImgUrl   string                      `json:"coverImgUrl"` //商品封面图
	Name          string                      `json:"name"`        // 活动名
	PriceYuan     string                      `json:"priceYuan"`
	LastSec       int64                       `json:"lastSec"`
	Status        int                         `json:"status"`
	Alias         string                      `json:"alias"`
	AssetIntro    string                      `json:"assetIntro"`
	ActivityIntro string                      `json:"activityIntro"`
	ActivityType  int                         `json:"activityType"`
	SubSumPeople  int                         `json:"subSumPeople"`
	SubSum        int                         `json:"subSum"`
	Subed         bool                        `json:"subed"`     //是否已认购
	Award         int                         `json:"award"`     //是否中签
	PayStatus     int                         `json:"payStatus"` //是否已付款
	Steps         []SubscribeActivityFullStep `json:"steps"`

	AssetCateString string      `json:"assetCateString"`
	AssetTotal      int         `json:"assetTotal"`
	AssetCreateAt   *gtime.Time `json:"assetCreateAt"`
	AssetDetailImg  string      `json:"assetDetailImg"`
	NfrDay          int         `json:"nfrDay"`
	ChainName       string      `json:"chainName"`
	AssetPic        string      `json:"assetPic"`
	ChainAddr       string      `json:"chainAddr"`
	ChainType       int         `json:"chainType"`

	CreatorId     int                         `orm:"creator_id" json:"creatorId"`
	CreatorName   string                      `orm:"creator_name" json:"creatorName"`
	CreatorAvatar string                      `orm:"creator_avatar" json:"creatorAvatar"`
	CreatorNo     string                      `orm:"creator_no" json:"creatorNo"`
	CopyrightInfo []provider.TplCopyrightInfo `json:"copyrightInfo"`
	AnHourAgo     int                         `json:"anHourAgo"` //活动结束前一个小时
}

type SubscribeActivityFullStep struct {
	Txt     string `json:"txt"`
	TimeStr string `json:"timeStr"`
}

type DoSubReq struct {
	SubNum             int    `json:"subNum"`
	Type               string `json:"type"`
	Alias              string `json:"alias"`
	UserId             string `json:"userId"`
	ClientIp           string `json:"clientIp"`
	OrderNo            string `json:"orderNo"`
	SuccessRedirectUrl string `json:"successRedirectUrl"`
	ExitRedirectUrl    string `json:"exitRedirectUrl"`
	PublisherId        string `json:"publisherId"`
	PlatformAppId      string `json:"platformAppId"`
}

type DoSubResult struct {
	Reason  string `json:"reason"`
	OrderNo string `json:"orderNo"`
	Type    string `json:"type"`
	Step    string `json:"step"`
}

const STATUS_AWAY_START = 1 //未开始
const STATUS_ING = 2        //进行中
const STATUS_AWAIT_OPEN = 3 //待公布
const STATUS_AWAIT_PAY = 4  //待付款
const STATUS_END = 5        //已结束
const STATUS_SUBED = 6      //已认购

const TICKET_MONTH = "month_ticket"
const TICKET_CRYSTAL = "crystal"
const TICKET_MONEY = "money"

const AWARD_STATUS_ING = 1 //0.未开奖；1.开奖中；2.开奖完毕
const AWARD_STATUS_END = 2

type AdminSubscribeActivityFull struct {
	SubscribeActivity
	Status    string `json:"status"`
	StatusTxt string `json:"statusTxt"`
	PriceYuan string `json:"priceYuan"`
}

type AdminListByPage struct {
	Total int                          `json:"total"`
	List  []AdminSubscribeActivityFull `json:"list"`
}

type AdminSubscribeActivityDetail struct {
	SubscribeActivity
	PriceYuan string               `json:"priceYuan"`
	Cons      []SubscribeCondition `json:"cons"`
}
