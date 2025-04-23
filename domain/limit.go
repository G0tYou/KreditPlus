package domain

import (
	"context"
	"database/sql"
)

type Limit struct {
	ID          int `json:"id"`
	LimitTypeID int `json:"limit_type_id" validate:"required"`
	UserID      int `json:"used_id" validate:"required"`
	Amount      int `json:"amount" validate:"required"`
	UsedAmount  int `json:"used_amount" validate:"required"`
}

type RepositoryMySQLLimit interface {
	//Listing repository
	// ReadTXLimitByID is a method to get limit by user_id and limit_type_id from the database
	ReadTXLimitByID(ctx context.Context, tx *sql.Tx, lid int) (Limit, error)

	//Updating repository
	// UpdateTXLimit is a method to update limit in the database
	UpdateTXLimit(ctx context.Context, tx *sql.Tx, l *Limit) error
}
