package repo

import (
	"github.com/yogenyslav/pkg/storage"
)

// AuthRepo implements methods for authorization on data layer.
type AuthRepo struct {
	pg storage.SQLDatabase
}

// New creates new AuthRepo.
func New(pg storage.SQLDatabase) *AuthRepo {
	return &AuthRepo{pg: pg}
}
