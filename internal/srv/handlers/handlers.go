package handlers

import (
	"gophermart-points/internal/repo/pgsql"

	"go.uber.org/zap"
)

const USER_ID_KEY = "userId"

type RsDef struct {
	Err string `json:"error"`
}

type API struct {
	db     *pgsql.PgSQL
	logger *zap.SugaredLogger
}

func NewAPI(db *pgsql.PgSQL, logger *zap.SugaredLogger) *API {
	return &API{
		db,
		logger,
	}
}
