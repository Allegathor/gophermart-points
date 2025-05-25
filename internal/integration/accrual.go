package integration

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type AccrualService struct {
	logger *zap.SugaredLogger
	addr   string
	Client *resty.Client
}

func NewInstance(addr string, logger *zap.SugaredLogger) *AccrualService {
	retryCount := 4

	c := resty.New()
	c.
		SetRetryCount(retryCount).
		SetRetryWaitTime(0).
		SetRetryAfter(func(c *resty.Client, r *resty.Response) (time.Duration, error) {
			attempt := r.Request.Attempt

			if attempt > retryCount {
				return 0, fmt.Errorf("max retries reached")
			}

			delay := time.Second + time.Duration(attempt-1)*2*time.Second
			logger.Infof("Retry attempt %d, waiting %v\n", attempt, delay)
			switch r.StatusCode() {
			case http.StatusNoContent:
				delay = 2*time.Second + time.Duration(attempt-1)*2*time.Second
			case http.StatusTooManyRequests:
				ra, err := strconv.Atoi(r.Header().Get("Retry-After"))
				if err != nil {
					logger.Errorln("failed atoi")
					return delay, nil
				}
				delay = time.Duration(ra) * time.Second
			}

			return delay, nil
		})

	s := &AccrualService{
		logger: logger,
		addr:   addr,
		Client: c,
	}

	return s
}

func (accr *AccrualService) UpdateOrderStatus(num string) (UpdateResult, error) {
	var res UpdateResult
	_, err := accr.Client.R().
		SetResult(&res).
		Get(accr.addr + "/api/orders/" + num)

	if err != nil {
		return res, err
	}
	accr.logger.Infow(
		"get order from accrual service",
		"response:", res,
	)

	return res, nil
}
