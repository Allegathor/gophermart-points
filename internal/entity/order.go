package entity

import (
	"time"
)

const (
	PointsEvalStatusNew        = "NEW"
	PointsEvalStatusProcessing = "PROCESSING"
	PointsEvalStatusProcessed  = "PROCESSED"
	PointsEvalStatusInvalid    = "INVALID"
)

type Order struct {
	UserID         int
	OrderID        int
	Num            string
	Amount         float64
	PntsEvalStatus string
	UploadAt       time.Time
}

func NewOrder(id int, num string, amount float64) *Order {
	return &Order{
		UserID:         id,
		Num:            num,
		Amount:         amount,
		PntsEvalStatus: PointsEvalStatusNew,
	}
}
