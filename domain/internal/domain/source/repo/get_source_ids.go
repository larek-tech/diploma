package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/pkg/errs"
)

const getSourceIDs = `
	select external_id
	from domain.source
	where internal_id in (
		select internal_source_id
		from domain.get_permitted_sources($2, $3)
		where internal_source_id = $1
	);
`

// GetSourceIDs returns external source id (uuid from data service).
func (r *Repo) GetSourceIDs(ctx context.Context, sourceID, userID int64, roleIDs []int64) (uuid.UUID, error) {
	var externalID uuid.UUID
	if err := r.pg.Query(ctx, &externalID, getSourceIDs, sourceID, userID, roleIDs); err != nil {
		return externalID, errs.WrapErr(err, "get source ids")
	}
	return externalID, nil
}
