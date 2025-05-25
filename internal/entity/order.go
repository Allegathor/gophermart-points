package entity

import (
	"time"
)

const (
	POINTS_EVAL_STATUS_NEW        = "NEW"
	POINTS_EVAL_STATUS_PROCESSING = "PROCESSING"
	POINTS_EVAL_STATUS_PROCESSED  = "PROCESSED"
	POINTS_EVAL_STATUS_INVALID    = "INVALID"
)

type Order struct {
	UserId         int
	OrderId        int
	Num            string
	Amount         float64
	PntsEvalStatus string
	UploadAt       time.Time
}

func NewOrder(id int, num string, amount float64) *Order {
	return &Order{
		UserId:         id,
		Num:            num,
		Amount:         amount,
		PntsEvalStatus: POINTS_EVAL_STATUS_NEW,
	}
}
