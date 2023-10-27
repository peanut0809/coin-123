package model

type CoinItems struct {
	Id              int    `json:"id"`
	CoinName        string `orm:"coin_name" json:"coinName"`
	DateContext     string `orm:"date_context" json:"dateContext"`
	Context         string `orm:"context" json:"context"`
	Type            int    `orm:"type" json:"type"`
	OriginCreatedAt string `orm:"origin_created_at" json:"originCreatedAt"`
	CreatedAt       string `orm:"created_at" json:"createdAt"`
	UpdatedAt       string `orm:"updated_at" json:"updatedAt"`
}
type CoinListReq struct {
	Page        int    `json:"pageNum"`
	PageSize    int    `json:"pageSize"`
	CoinName    string `json:"CoinName"`
	Type        int    `json:"type"`
	DateContext string `json:"dateContext"`
}

type CoinList struct {
	List  []*CoinItems `json:"list"`
	Total int          `json:"total"`
}
