package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func IsNoRowsError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
