package model

import "time"

type SyntheticRecord struct {
	Id        int       `orm:"id" json:"id"`
	OrderNo   string    `orm:"order_no" json:"orderNo"`     // 订单号
	AssetName string    `orm:"asset_name" json:"assetName"` // 资产名
	AssetIcon string    `orm:"asset_icon" json:"assetIcon"` // 图标
	AssetPic  string    `orm:"asset_pic" json:"assetPic"`   // 资产封面
	Aid       int       `orm:"aid" json:"aid"`              // 活动ID
	InData    string    `orm:"in_data" json:"inData"`       // 合成数据
	OutData   string    `orm:"out_data" json:"outData"`     // 产出数据
	CreatedAt time.Time `orm:"created_at" json:"createdAt"` // 新建时间
	UpdatedAt time.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}
