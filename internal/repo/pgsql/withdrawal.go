package pgsql

import (
	"context"
	"gophermart-points/internal/entity"
	"math"

	"github.com/jackc/pgx/v5"
)

func GetWithdrawalSumQuery(pg PgDB, ctx context.Context, userID int) (sum float64, err error) {
	row := pg.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0) FROM transaction_table
		WHERE user_id = @id
		AND amount < 0
	`,
		pgx.NamedArgs{"id": userID},
	)

	err = row.Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (pg *PgSQL) GetWithdrawalSum(ctx context.Context, userID int) (sum float64, err error) {
	err = pg.WithPolicy(
		ctx,
		func(ctx context.Context) error {
			sum, err = GetWithdrawalSumQuery(pg, ctx, userID)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return 0, err
	}

	return math.Abs(sum), nil
}

func (pg *PgSQL) GetWithdrawals(ctx context.Context, userID int) ([]entity.Withdrawal, error) {
	wls := make([]entity.Withdrawal, 0)
	err := pg.WithPolicy(
		ctx,
		func(c context.Context) error {
			rows, err := pg.Query(c, `
				SELECT tt.user_id, ot.order_num, tt.amount, tt.processed_at  
				FROM transaction_table AS tt
				JOIN order_table AS ot 
					ON tt.order_id = ot.order_id 
					AND tt.user_id = @id
					AND tt.amount < 0
				ORDER BY processed_at DESC;
			`,
				pgx.NamedArgs{"id": userID},
			)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var w entity.Withdrawal
				err := rows.Scan(&w.UserID, &w.Num, &w.Amount, &w.ProcAt)
				if err != nil {
					return err
				}
				wls = append(wls, w)
			}

			return nil
		},
	)
	if err != nil {
		return wls, err
	}

	return wls, nil
}
