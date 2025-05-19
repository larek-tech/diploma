package handler

import (
	"github.com/larek-tech/diploma/api/internal/domain/pb"
)

const (
	sourceIDParam = "id"
	domainIDParam = "id"
	offsetParam   = "offset"
	limitParam    = "limit"
)

// Handler implements source methods on transport level.
type Handler struct {
	sourceService pb.SourceServiceClient
}

// New creates new Handler.
func New(sourceService pb.SourceServiceClient) *Handler {
	return &Handler{
		sourceService: sourceService,
	}
}
