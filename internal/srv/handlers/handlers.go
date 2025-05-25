package handlers

import "gophermart-points/internal/repo/pgsql"

const USER_ID_KEY = "userId"

type RsDef struct {
	Err string `json:"error"`
}

type API struct {
	db *pgsql.PgSQL
}

func NewAPI(db *pgsql.PgSQL) *API {
	return &API{
		db,
	}
}
