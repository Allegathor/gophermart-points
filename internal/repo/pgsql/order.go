package pgsql

import (
	"context"
	"errors"
	"fmt"
	"gophermart-points/internal/entity"

	"github.com/jackc/pgx/v5"
)

func AddOrderQuery(pg PgDB, ctx context.Context, order entity.Order) (orderId int, err error) {
	row := pg.QueryRow(ctx, `
		INSERT INTO order_table (user_id, order_num, amount, points_eval_status)
		VALUES(@id, @num, @amount, @status)
		RETURNING order_id;
	`,
		pgx.NamedArgs{
			"id":     order.UserID,
			"num":    order.Num,
			"amount": order.Amount,
			"status": order.PntsEvalStatus,
		},
	)

	err = row.Scan(&orderId)
	if err != nil {
		return -1, err
	}

	return orderId, nil
}

func UpdateEvalPntsStatusQuery(pg PgDB, ctx context.Context, userID int, orderId int, status string) error {
	_, err := pg.Exec(ctx, `
		UPDATE order_table	
		SET points_eval_status = @s
		WHERE user_id = @userID AND order_id = @orderId
		RETURNING order_id;
	`,
		pgx.NamedArgs{
			"userID":  userID,
			"orderId": orderId,
			"s":       status,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEvalPntsQuery(pg PgDB, ctx context.Context, order entity.Order) error {
	_, err := pg.Exec(ctx, `
		UPDATE order_table	
		SET amount = @amount, points_eval_status = @s
		WHERE user_id = @userID AND order_id = @orderId
		RETURNING order_id;
	`,
		pgx.NamedArgs{
			"userID":  order.UserID,
			"orderId": order.OrderId,
			"amount":  order.Amount,
			"s":       order.PntsEvalStatus,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PgSQL) AddOrder(ctx context.Context, order entity.Order) (orderId int, err error) {
	err = pg.WithPolicy(
		ctx,
		func(ctx context.Context) error {
			orderId, err = AddOrderQuery(pg, ctx, order)
			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return -1, err
	}

	return orderId, nil
}

func (pg *PgSQL) GetOrder(ctx context.Context, orderNum string) (entity.Order, error) {
	var order entity.Order
	err := pg.WithPolicy(
		ctx,
		func(c context.Context) error {
			row := pg.QueryRow(c, `
				SELECT user_id, order_id, order_num, amount, points_eval_status  
				FROM order_table
				WHERE order_num = @num
			`,
				pgx.NamedArgs{"num": orderNum},
			)

			err := row.Scan(&order.UserID, &order.OrderId, &order.Num, &order.Amount, &order.PntsEvalStatus)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return order, ErrNoOrder
		}

		return order, err
	}

	return order, nil
}

func (pg *PgSQL) GetOrders(ctx context.Context, userID int) ([]entity.Order, error) {
	orders := make([]entity.Order, 0)
	err := pg.WithPolicy(
		ctx,
		func(c context.Context) error {
			rows, err := pg.Query(c, `
				SELECT user_id, order_num, amount, points_eval_status, uploaded_at  
				FROM order_table
				WHERE user_id = @id
				AND amount >= 0
				ORDER BY uploaded_at DESC
			`,
				pgx.NamedArgs{"id": userID},
			)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var o entity.Order
				err := rows.Scan(&o.UserID, &o.Num, &o.Amount, &o.PntsEvalStatus, &o.UploadAt)
				if err != nil {
					return err
				}
				fmt.Println(o)
				orders = append(orders, o)
			}

			return nil
		},
	)
	if err != nil {
		return orders, err
	}

	fmt.Println(orders)
	return orders, nil
}

func (pg *PgSQL) UpdateEvalPntsStatus(ctx context.Context, userID int, orderId int, status string) error {
	return pg.WithPolicy(
		ctx,
		func(ctx context.Context) error {
			return UpdateEvalPntsStatusQuery(pg, ctx, userID, orderId, status)
		},
	)
}

func (pg *PgSQL) UpdateEvalPnts(ctx context.Context, order entity.Order) error {
	return pg.WithPolicy(
		ctx,
		func(ctx context.Context) error {
			return UpdateEvalPntsQuery(pg, ctx, order)
		},
	)
}
