package entity

type Transaction struct {
	UserID   int    `json:"userID,omitempty"`
	OrderID  int    `json:"orderId,omitempty"`
	OrderNum string `json:"num"`
	Amount   int    `json:"sum"`
}
