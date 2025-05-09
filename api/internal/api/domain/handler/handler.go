package handler

import (
	"github.com/larek-tech/diploma/api/internal/domain/pb"
)

const (
	domainIDParam = "id"
	offsetParam   = "offset"
	limitParam    = "limit"
)

// Handler implements domain methods on transport level.
type Handler struct {
	domainService pb.DomainServiceClient
}

// New creates new Handler.
func New(domainService pb.DomainServiceClient) *Handler {
	return &Handler{
		domainService: domainService,
	}
}
