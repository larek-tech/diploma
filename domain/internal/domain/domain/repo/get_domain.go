package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/domain/model"
	"github.com/yogenyslav/pkg/errs"
)

const getDomainByID = `
	select id, title, user_id, source_ids, created_at, updated_at
	from domain.domain
		where id = (
			select id
			from domain.get_permitted_domains($2, $3)
			where id = $1
		);
`

// GetDomainByID returns domain by ID.
func (r *Repo) GetDomainByID(ctx context.Context, id, userID int64, roleIDs []int64) (model.DomainDao, error) {
	var domain model.DomainDao
	if err := r.pg.Query(ctx, &domain, getDomainByID, id, userID, roleIDs); err != nil {
		return domain, errs.WrapErr(err, "get domain by id")
	}
	return domain, nil
}
