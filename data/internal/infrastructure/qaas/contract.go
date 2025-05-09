package qaas

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas/messages"
)

type (
	pageJobHandler interface {
		Handle(ctx context.Context, job messages.PageJob) error
	}
	siteJobHandler interface {
		Handle(ctx context.Context, job messages.SiteJob) error
	}
	resultMessageHandler interface {
		Handle(ctx context.Context, result messages.ResultMessage) error
	}
)
