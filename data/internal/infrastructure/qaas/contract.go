package qaas

import (
	"context"

	"go.dataddo.com/pgq"
)

type (
	handler interface {
		Handle(context.Context, *pgq.MessageIncoming) (bool, error)
	}
)
	