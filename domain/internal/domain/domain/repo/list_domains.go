package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/domain/model"
	"github.com/yogenyslav/pkg/errs"
)

const listDomains = `
	select id, title, user_id, source_ids, created_at, updated_at
	from domain.domain
		where id in (
			select id
			from domain.get_permitted_domains($1, $2)
		)
	order by created_at desc, updated_at desc
	offset $3
	limit $4;
`

// ListDomains returns list of domains available for user.
func (r *Repo) ListDomains(ctx context.Context, userID int64, roleIDs []int64, offset, limit uint64) ([]model.DomainDao, error) {
	var domains []model.DomainDao
	if err := r.pg.QuerySlice(ctx, &domains, listDomains, userID, roleIDs, offset, limit); err != nil {
		return domains, errs.WrapErr(err, "list domains")
	}
	return domains, nil
}
