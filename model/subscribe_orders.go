package model

type SubscribeOrder struct {
	Id      int    `orm:"id" json:"id"`
	OrderNo string `orm:"order_no" json:"orderNo"`
	Price   int    `orm:"price" json:"price"`
	Status  int    `orm:"status" json:"status"`
	UserId  string `orm:"user_id" json:"userId"`
}
