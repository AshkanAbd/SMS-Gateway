package models

import "errors"

var (
	EmptyNameError        = errors.New("user name is empty")
	UserNotExistError     = errors.New("user does not exist")
	NotEnoughBalanceError = errors.New("not enough balance")
)
