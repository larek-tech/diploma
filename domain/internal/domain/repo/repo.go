package repo

import (
	"github.com/yogenyslav/pkg/storage"
)

// DomainRepo implements domain methods on data layer.
type DomainRepo struct {
	pg storage.SQLDatabase
}

// NewDomainRepo creates new DomainRepo.
func NewDomainRepo(pg storage.SQLDatabase) *DomainRepo {
	return &DomainRepo{
		pg: pg,
	}
}

// SourceRepo implements source methods on data layer.
type SourceRepo struct {
	pg storage.SQLDatabase
}

// NewSourceRepo creates new SourceRepo.
func NewSourceRepo(pg storage.SQLDatabase) *SourceRepo {
	return &SourceRepo{
		pg: pg,
	}
}
