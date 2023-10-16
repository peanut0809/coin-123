package model

type TransactionSlip struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Trade struct {
	Count     int    `json:"count"`
	CreatedAt string `json:"createdAt"`
}

type Price struct {
	Price     float64 `json:"price"`
	CreatedAt string  `json:"createdAt"`
}
