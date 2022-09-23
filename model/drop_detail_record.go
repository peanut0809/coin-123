package model

import (
	"github.com/gogf/gf/os/gtime"
)

type DropDetailRecord struct {
	Id         int         `orm:"id" json:"id"`            // 主键
	Phone      string      `orm:"phone" json:"phone"`      // 手机号
	UserId     string      `orm:"user_id" json:"userId"`   // 用户ID
	DropId     int         `orm:"drop_id" json:"dropId"`   // 空投ID
	Status     int         `orm:"status" json:"status"`    // 0.待空投；1.空投成功；2.空投失败；
	AppId      string      `orm:"app_id" json:"appId"`     // appid
	TokenId    string      `orm:"token_id" json:"tokenId"` // tokenid
	TemplateId string      `orm:"template_id" json:"templateId"`
	ErrMsg     string      `orm:"err_msg" json:"errMsg"` // 错误信息
	NfrSec     int         `orm:"nfr_sec" json:"nfrSec"`
	CreatedAt  *gtime.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt  *gtime.Time `orm:"updated_at" json:"updatedAt"` // 新建时间
}

type DropDetailRecordList struct {
	Total int                `json:"total"`
	List  []DropDetailRecord `json:"list"`
}
