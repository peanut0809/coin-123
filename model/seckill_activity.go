package model

import (
	"github.com/gogf/gf/os/gtime"
)

type SeckillActivity struct {
	Id                int         `orm:"id,primary" json:"id"` // 活动iD
	Alias             string      `orm:"alias" json:"alias"`   // 别名
	Disable           int         `orm:"disable" json:"disable"`
	NfrSec            int         `orm:"nfr_sec" json:"nfrSec"`
	PublisherId       string      `orm:"publisher_id" json:"publisherId"`
	ActivityStartTime *gtime.Time `orm:"activity_start_time" json:"activityStartTime"` // 活动开始时间
	ActivityEndTime   *gtime.Time `orm:"activity_end_time" json:"activityEndTime"`     //活动结束时间
	StartTime         *gtime.Time `orm:"start_time" json:"startTime"`                  // 普通用户认购开始时间
	Price             int         `orm:"price" json:"price"`                           // 发售价,单位：分
	AppId             string      `orm:"app_id" json:"appId"`                          // 应用ID
	AssetType         string      `orm:"asset_type" json:"assetType"`                  // 资产类型
	AssetType2        string      `orm:"asset_type2" json:"assetType2"`                // 资产类型
	AssetType3        string      `orm:"asset_type3" json:"assetType3"`                // 资产类型
	TemplateId        string      `orm:"template_id" json:"templateId"`                // 模板ID
	SumNum            int         `orm:"sum_num" json:"sumNum"`                        // 总数
	RemainNum         int         `orm:"remain_num" json:"remainNum"`
	LimitBuy          int         `orm:"limit_buy" json:"limitBuy"`
	ActivityIntro     string      `orm:"activity_intro" json:"activityIntro"`
	CoverImgUrl       string      `orm:"cover_img_url" json:"coverImgUrl"`
	Name              string      `orm:"name" json:"name"`
	CreatedAt         *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt         *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}

type SeckillActivityFull struct {
	*SeckillActivity
	LastSec   int64  `json:"lastSec"`
	Status    int    `json:"status"`
	PriceYuan string `json:"priceYuan"`

	AssetCateString string      `json:"assetCateString"`
	AssetTotal      int         `json:"assetTotal"`
	AssetCreateAt   *gtime.Time `json:"assetCreateAt"`
	AssetDetailImg  string      `json:"assetDetailImg"`
	AssetPic        string      `json:"assetPic"`
	NfrDay          int         `json:"nfrDay"`
	ChainName       string      `json:"chainName"`
	ChainAddr       string      `json:"chainAddr"`
	ChainType       int         `json:"chainType"`
}

const SeckillActivityStatus_Wait_Start = 0
const SeckillActivityStatus_Ing = 1
const SeckillActivityStatus_End = 2

type DoBuyReq struct {
	Alias              string `json:"alias"`
	Num                int    `json:"num"`
	ClientIp           string `json:"clientIp"`
	SuccessRedirectUrl string `json:"successRedirectUrl"`
	ExitRedirectUrl    string `json:"exitRedirectUrl"`
	PublisherId        string `json:"publisherId"`
	PlatformAppId      string `json:"platformAppId"`
	UserId             string `json:"userId"`
	OrderNo            string `json:"orderNo"`
}

type CreateSeckillActivityReq struct {
	SeckillActivity
	PriceYuan string `json:"priceYuan"`
}

type AdminSeckillActivityList struct {
	List  []AdminSeckillActivityFull `json:"list"`
	Total int                        `json:"total"`
}

type AdminSeckillActivityFull struct {
	SeckillActivity
	PriceYuan string `json:"priceYuan"`
	StatusTxt string `json:"statusTxt"`
	Status    string `json:"status"`
}
