package pgsql

import (
	"context"
	"gophermart-points/internal/entity"

	"github.com/jackc/pgx/v5"
)

func GetBalanceQuery(pg PgDB, ctx context.Context, userId int) (sum float64, err error) {
	row := pg.QueryRow(ctx, `
		SELECT * FROM get_current_user_balance(@id);
	`,
		pgx.NamedArgs{"id": userId},
	)

	err = row.Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (pg *PgSQL) GetBalance(ctx context.Context, userId int) (sum float64, err error) {
	err = pg.WithPolicy(
		ctx,
		func(ctx context.Context) error {
			sum, err = GetBalanceQuery(pg, ctx, userId)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (pg *PgSQL) Accrue(ctx context.Context, order entity.Order) error {
	return pg.TxWithPolicy(
		ctx,
		pgx.TxOptions{AccessMode: pgx.ReadWrite},
		func(tx pgx.Tx, c context.Context) error {

			order.PntsEvalStatus = entity.POINTS_EVAL_STATUS_PROCESSED
			err := UpdateEvalPntsQuery(tx, c, order)
			if err != nil {
				return err
			}

			err = AddTransactionQuery(tx, c, order)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func (pg *PgSQL) Withdraw(ctx context.Context, wd entity.Withdrawal) error {
	return pg.TxWithPolicy(
		ctx,
		pgx.TxOptions{AccessMode: pgx.ReadWrite},
		func(tx pgx.Tx, c context.Context) error {
			balance, err := GetBalanceQuery(tx, c, wd.UserId)
			if err != nil {
				return err
			} else if wd.AbsAmount > balance {
				return ErrInsufficentFunds
			}

			order := entity.NewOrder(wd.UserId, wd.Num, wd.Amount)
			orderId, err := AddOrderQuery(tx, c, *order)
			if err != nil {
				return err
			}
			order.OrderId = orderId

			err = AddTransactionQuery(tx, c, *order)
			if err != nil {
				return err
			}

			err = UpdateEvalPntsStatusQuery(tx, c, wd.UserId, order.OrderId, entity.POINTS_EVAL_STATUS_PROCESSED)
			if err != nil {
				return err
			}

			return nil
		},
	)

}
