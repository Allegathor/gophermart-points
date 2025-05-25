package pgsql

import (
	"context"
	"errors"
	"gophermart-points/internal/entity"

	"github.com/jackc/pgx/v5"
)

func (pg *PgSQL) HasUser(ctx context.Context, login string) (bool, error) {
	err := pg.WithPolicy(
		ctx,
		func(c context.Context) error {
			var val string
			row := pg.QueryRow(c, `
				SELECT login
				FROM user_table
				WHERE login = @l
			`,
				pgx.NamedArgs{"l": login},
			)

			err := row.Scan(&val)
			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (pg *PgSQL) CreateUser(ctx context.Context, login, passwd string) (id int, err error) {
	err = pg.WithPolicy(
		ctx,
		func(c context.Context) error {
			err := pg.QueryRow(c, `
				INSERT INTO user_table (login, passwd)
				VALUES(@login, @passwd)
				RETURNING user_id;
			`,
				pgx.NamedArgs{"login": login, "passwd": passwd},
			).Scan(&id)

			return err
		},
	)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (pg *PgSQL) GetUser(ctx context.Context, login string) (u entity.User, err error) {
	err = pg.WithPolicy(
		ctx,
		func(c context.Context) error {
			err := pg.QueryRow(c, `
				SELECT user_id, login, passwd FROM user_table
				WHERE login = @l;
			`,
				pgx.NamedArgs{"l": login},
			).Scan(&u.ID, &u.Login, &u.Pwd)

			return err
		},
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, ErrUnexistLogin
		}

		return u, err
	}

	return u, nil
}
