package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("your requested Item is not found")
	ErrConflict            = errors.New("your Item already exist")
	ErrConflictUsername    = errors.New("your Username already exist")
	ErrBadParamInput       = errors.New("given Param is not valid")
	ErrUserNotFound        = errors.New("your Username is not found")
	ErrLogin               = errors.New("login failed, wrong username or password")
	ErrGenerateToken       = errors.New("error while generating token")
	ErrLimitExceeded       = errors.New("limit amount exceeded")
)
