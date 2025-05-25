package pgsql

import "errors"

var (
	ErrInsufficentFunds = errors.New("insufficient funds")
	ErrUnexistLogin     = errors.New("this login is unexist")
	ErrNoOrder          = errors.New("there is no order with provided num")
)
