package model

import (
	"github.com/gogf/gf/os/gtime"
)

type DropRecord struct {
	Id         int         `orm:"id" json:"id"`     // 主键
	Name       string      `orm:"name" json:"name"` // 名字
	Phones     string      `orm:"phones" json:"phones"`
	OrderNo    string      `orm:"order_no" json:"orderNo"`       // 订单号
	AppId      string      `orm:"app_id" json:"appId"`           // appid
	TemplateId string      `orm:"template_id" json:"templateId"` // 模板ID
	Num        int         `orm:"num" json:"num"`                // 空投数量
	NfrSec     int         `orm:"nfr_sec" json:"nfrSec"`
	Remark     string      `orm:"remark" json:"remark"`
	CreatedAt  *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt  *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}

type DropRecordReq struct {
	DropRecord
	PhoneArr  []string `json:"phoneArr"`
	ExcelFile string   `json:"excelFile"`
}

type DropRecordFull struct {
	DropRecord
	AssetName string `json:"assetName"`
}

type RecordList struct {
	Total int              `json:"total"`
	List  []DropRecordFull `json:"list"`
}
