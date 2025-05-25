package handlers

import (
	"context"
	"gophermart-points/internal/entity"

	"github.com/golodash/galidator/v2"
)

const USER_COOKIE_NAME = "gmp_auth"

var (
	gu = galidator.New()
	// userValid = gu.ComplexValidator(galidator.Rules{
	// 	"Login": gu.R("login").Regex("[A-Za-z0-9_]").Required().Min(3).Max(32),
	// 	"Pwd":   gu.R("password").Required().Password(),
	// })
	userValid = gu.ComplexValidator(galidator.Rules{
		"Login": gu.R("login").Required().Min(3).Max(32),
		"Pwd":   gu.R("password").Required(),
	})
)

type RsUser struct {
	Err       string `json:"error"`
	FieldErrs any    `json:"fieldErrors"`
}

type UserRepo interface {
	HasUser(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, login, pwd string) (id int, err error)
	GetUser(ctx context.Context, login string) (user entity.User, err error)
}

type UserAPI struct {
	authKey string
	db      UserRepo
}

func NewUserAPI(authKey string, db UserRepo) *UserAPI {
	return &UserAPI{
		authKey,
		db,
	}
}
