package integration

const (
	StatusRegistered = "REGISTERED"
	StatusInvalid    = "INVALID"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
)

type UpdateResult struct {
	Num    string  `json:"order"`
	Status string  `json:"status"`
	Amount float64 `json:"accrual"`
}
