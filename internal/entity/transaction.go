package entity

type Transaction struct {
	userID   int    `json:"userID,omitempty"`
	OrderId  int    `json:"orderId,omitempty"`
	OrderNum string `json:"num"`
	Amount   int    `json:"sum"`
}
