package entity

import "time"

type Withdrawal struct {
	UserId    int       `json:"userId,omitempty"`
	Num       string    `json:"order"`
	AbsAmount float64   `json:"sum"`
	Amount    float64   `json:"oppAmount,omitempty"`
	ProcAt    time.Time `json:"procAt"`
}
