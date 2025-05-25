package pgsql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PgDB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type PgSQL struct {
	*pgxpool.Pool
	logger *zap.SugaredLogger
}

func (pg *PgSQL) Close() {
	pg.Pool.Close()
}

func IsRetryable(err error) bool {
	var pgErr *pgconn.PgError
	if err == nil {
		return false
	}

	if err == context.DeadlineExceeded || err == context.Canceled {
		return true
	}

	if !errors.As(err, &pgErr) {
		return false
	}

	switch pgErr.Code {
	case
		pgerrcode.SerializationFailure,
		pgerrcode.DeadlockDetected:
		return true
	}

	return pgerrcode.IsConnectionException(pgErr.Code)
}

func calcDelay(exec failsafe.ExecutionAttempt[any]) time.Duration {
	return time.Second + time.Duration(exec.Attempts()-1)*2*time.Second
}

const maxRetryes = 3

var retry = retrypolicy.Builder[any]().HandleIf(func(_ any, err error) bool {
	return IsRetryable(err)
}).
	WithMaxRetries(maxRetryes).
	WithDelayFunc(calcDelay).
	Build()

func (pg *PgSQL) WithPolicy(
	ctx context.Context, fn func(context.Context) error,
) error {

	return failsafe.NewExecutor(retry).
		WithContext(ctx).
		RunWithExecution(func(exec failsafe.Execution[any]) (err error) {
			ectx := exec.Context()
			if err := fn(ectx); err != nil {
				return err
			}

			return nil
		})
}

func (pg *PgSQL) TxWithPolicy(
	ctx context.Context, txOptions pgx.TxOptions, fn func(pgx.Tx, context.Context) error,
) error {

	return failsafe.NewExecutor(retry).
		WithContext(ctx).
		RunWithExecution(func(exec failsafe.Execution[any]) (err error) {
			ectx := exec.Context()
			tx, err := pg.BeginTx(ectx, txOptions)
			if err != nil {
				return fmt.Errorf("failed to begin tx: %w", err)
			}

			defer tx.Rollback(ectx)

			if err := fn(tx, ectx); err != nil {
				return err
			}

			if err := tx.Commit(ectx); err != nil {
				return fmt.Errorf("failed to commit: %w", err)
			}

			return err
		})
}

func Init(ctx context.Context, connStr string, logger *zap.SugaredLogger) (*PgSQL, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 15
	config.MinConns = 2
	config.MaxConnIdleTime = 20 * time.Second
	config.HealthCheckPeriod = 10 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	pg := &PgSQL{Pool: pool, logger: logger}

	return pg, nil
}
