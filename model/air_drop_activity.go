package model

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

const Air_Drop_Status_01 = 1
const Air_Drop_Status_02 = 2

const Air_Drop_Type_Speed = "speed"
const Air_Drop_Type_Crystal = "crystal"

var Air_Drop_Type_Map = g.Map{
	Air_Drop_Type_Speed:   "用户加速次数空投",
	Air_Drop_Type_Crystal: "用户元晶空投",
}

type AirDropActivity struct {
	Id           int         `orm:"id" json:"id"`                // 主键
	Name         string      `orm:"name" json:"name"`            // 名字
	OrderNo      string      `orm:"order_no" json:"orderNo"`     // 订单编号
	Remark       string      `orm:"remark" json:"remark"`        // 备注
	Config       string      `orm:"config" json:"config"`        // 请求配置
	Type         string      `orm:"type" json:"type"`            // 空投活动类型 crystal-元晶 speed 加速次数
	CreatedAt    *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt    *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
	SuccessCount int         `json:"successCount"`               // 成功数量
	ErrorCount   int         `json:"errorCount"`                 // 失败数量
}

type AirDropActivityReq struct {
	DropId    int    `json:"dropId"`
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	Type      string `json:"type"`
	ExcelFile string `json:"excelFile"`
}

type AirDropActivityItemReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	DropId   int    `json:"dropId"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
}
type AirDropActivityItemsRex struct {
	List  []*AirDropActivity `json:"list"`
	Total int                `json:"total"`
}

type AirDropActivityItemRex struct {
	*AirDropActivity
	Items []*AirDropActivityRecord `json:"items"`
	Total int                      `json:"total"`
}

type MobileCollect struct {
	Mobile     string `json:"mobile"`
	Number     int    `json:"number"`
	UserId     string `json:"userId"`
	Message    string `json:"message"`
	HaveErr    bool   `json:"haveErr"`
	ActivityId int    `json:"activityId"`
}
type MarktingUserSpeedNumReq struct {
	UserId string `json:"userId"`
	Number int    `json:"number"`
	From   int    `json:"From"`
}
