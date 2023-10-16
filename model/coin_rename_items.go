package model

type CoinRenameItems struct {
	Id              int    `json:"id"`
	OriginalName    string `orm:"original_name" json:"originalName"`
	AfterName       string `orm:"after_name" json:"afterName"`
	Context         string `orm:"context" json:"context"`
	OriginCreatedAt string `orm:"origin_created_at" json:"originCreatedAt"`
	CreatedAt       string `orm:"created_at" json:"createdAt"`
	UpdatedAt       string `orm:"updated_at" json:"updatedAt"`
}
type CoinRenameReq struct {
	Page         int    `json:"pageNum"`
	PageSize     int    `json:"pageSize"`
	OriginalName string `json:"originalName"`
	AfterName    string `json:"afterName"`
}

type CoinRenameList struct {
	List  []*CoinRenameItems `json:"list"`
	Total int                `json:"total"`
}
