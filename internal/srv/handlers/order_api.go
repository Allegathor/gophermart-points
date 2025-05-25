package handlers

import (
	"gophermart-points/internal/repo/pgsql"
	"gophermart-points/internal/srv/external"

	"github.com/golodash/galidator/v2"
)

var gord = galidator.New()
var orderNumValid = gord.Validator(gord.R("ordernum").Regex("[0-9]").Required().Min(3).Max(32))

type OrderAPI struct {
	db    *pgsql.PgSQL
	Queue *external.OrderProcessing
}

func NewOrderAPI(db *pgsql.PgSQL, q *external.OrderProcessing) *OrderAPI {
	return &OrderAPI{
		db:    db,
		Queue: q,
	}
}

type RsOrder struct {
	Err       string `json:"error"`
	FieldErrs any    `json:"fieldErrors"`
}

type OrderRec struct {
	Num      string  `json:"number"`
	Status   string  `json:"status"`
	Amount   float64 `json:"accrual"`
	UploadAt string  `json:"uploaded_at"`
}
