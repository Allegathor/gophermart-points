package integration

const (
	STATUS_REGISTERED = "REGISTERED"
	STATUS_INVALID    = "INVALID"
	STATUS_PROCESSING = "PROCESSING"
	STATUS_PROCESSED  = "PROCESSED"
)

type UpdateResult struct {
	Num    string  `json:"order"`
	Status string  `json:"status"`
	Amount float64 `json:"accrual"`
}
