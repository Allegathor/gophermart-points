package entity

type Transaction struct {
	UserId   int    `json:"userId,omitempty"`
	OrderId  int    `json:"orderId,omitempty"`
	OrderNum string `json:"num"`
	Amount   int    `json:"sum"`
}
