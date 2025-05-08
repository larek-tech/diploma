package repo

import (
	"github.com/yogenyslav/pkg/storage"
)

// Repo implements chat methods on data layer.
type Repo struct {
	pg storage.SQLDatabase
}

// New creates new Repo.
func New(pg storage.SQLDatabase) *Repo {
	return &Repo{pg: pg}
}
