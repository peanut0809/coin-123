package model

type CoinIeoOffItems struct {
	Id              int    `json:"id"`
	Context         string `orm:"context" json:"context"`
	CoinName        string `orm:"coin_name" json:"coin_name"`
	Type            string `orm:"type" json:"type"`
	OriginCreatedAt string `orm:"origin_created_at" json:"originCreatedAt"`
	CreatedAt       string `orm:"created_at" json:"createdAt"`
	UpdatedAt       string `orm:"updated_at" json:"updatedAt"`
}
type CoinIeoOffReq struct {
	Page     int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Type     string `json:"type"`
	CoinName string `json:"coinName"`
}

type CoinIeoOffList struct {
	List  []*CoinIeoOffItems `json:"list"`
	Total int                `json:"total"`
}
