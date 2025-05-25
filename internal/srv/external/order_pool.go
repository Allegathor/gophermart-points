package external

import (
	"context"
	"gophermart-points/internal/entity"
	"gophermart-points/internal/integration"
	"gophermart-points/internal/repo/pgsql"
	"time"

	"go.uber.org/zap"
)

type OrderProcessing struct {
	ctx     context.Context
	db      *pgsql.PgSQL
	logger  *zap.SugaredLogger
	Accrual integration.AccrualService
	Chan    chan entity.Order
}

func NewOrderProcessing(ctx context.Context, db *pgsql.PgSQL, accrualAddr string, chanCap uint, logger *zap.SugaredLogger) *OrderProcessing {
	return &OrderProcessing{
		ctx:     ctx,
		db:      db,
		logger:  logger,
		Accrual: *integration.NewInstance(accrualAddr, logger),
		Chan:    make(chan entity.Order, chanCap),
	}
}

func (op *OrderProcessing) Push(o entity.Order) {
	op.Chan <- o
}

func (op *OrderProcessing) Process(idx uint) {
	for order := range op.Chan {
		ticker := time.NewTicker(time.Second * 2)
	outer:
		for {
			select {
			case <-ticker.C:
				{

					res, err := op.Accrual.UpdateOrderStatus(order.Num)
					if err != nil {
						op.logger.Errorln(err)
						err := op.db.UpdateEvalPntsStatus(op.ctx, order.UserID, order.OrderId, entity.PointsEvalStatusInvalid)
						if err != nil {
							op.logger.Errorln(err)
						}
						ticker.Stop()
						break outer
					}
					switch res.Status {
					case integration.StatusRegistered:
						op.logger.Infow("nothing to update from accrual", "order", order, "res", res)
					case integration.StatusProcessing:
						err := op.db.UpdateEvalPntsStatus(op.ctx, order.UserID, order.OrderId, entity.PointsEvalStatusProcessed)
						if err != nil {
							op.logger.Errorw(err.Error(), "order", order, "res", res)
							ticker.Stop()
							break outer
						}
						op.logger.Infow("update status from accrual service to PROCESSED", "order", order, "res", res)
					case integration.StatusInvalid:
						err := op.db.UpdateEvalPntsStatus(op.ctx, order.UserID, order.OrderId, entity.PointsEvalStatusInvalid)
						if err != nil {
							op.logger.Errorw(err.Error(), "order", order, "res", res)
							ticker.Stop()
							break outer
						}
						op.logger.Infow("update status from accrual service to INVALID", "order", order, "res", res)
						ticker.Stop()
						break outer
					case integration.StatusProcessed:
						order.Amount = res.Amount
						err := op.db.Accrue(op.ctx, order)
						if err != nil {
							op.logger.Errorw(err.Error(), "order", order, "res", res)
							ticker.Stop()
							break outer
						}
						op.logger.Infow("succesfull accrual", "order", order, "res", res)
						ticker.Stop()
						break outer
					}
				}
			case <-op.ctx.Done():
				{
					ticker.Stop()
				}

			}

		}
	}
}

func (op *OrderProcessing) RunPool(workerCount uint) {
	for i := range workerCount {
		go op.Process(i)
	}
}
