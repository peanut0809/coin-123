package model

type TransactionSlip struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Trade struct {
	Count     string `json:"count"`
	CreatedAt string `json:"createdAt"`
}
