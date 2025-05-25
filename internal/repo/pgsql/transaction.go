package pgsql

import (
	"context"
	"gophermart-points/internal/entity"

	"github.com/jackc/pgx/v5"
)

func AddTransactionQuery(pg PgDB, ctx context.Context, order entity.Order) error {
	_, err := pg.Exec(ctx, `
		INSERT INTO transaction_table (user_id, order_id, amount)
		VALUES(@id, @orderID, @amount)
	`,
		pgx.NamedArgs{
			"id":      order.UserID,
			"orderID": order.OrderID,
			"amount":  order.Amount,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
