package model

type CoinCoreContext struct {
	Id              int    `json:"id"`
	Context         string `orm:"context" json:"context"`
	OriginCreatedAt string `orm:"origin_created_at" json:"originCreatedAt"`
	CreatedAt       string `orm:"created_at" json:"createdAt"`
	UpdatedAt       string `orm:"updated_at" json:"updatedAt"`
}
type CoinCoreReq struct {
	Page     int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type CoinCoreList struct {
	List  []*CoinCoreContext `json:"list"`
	Total int                `json:"total"`
}
