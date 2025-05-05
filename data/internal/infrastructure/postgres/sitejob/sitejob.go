package sitejob

import "context"

type SiteJobStorage struct {
	db db
}

func New(db db) *SiteJobStorage {
	return &SiteJobStorage{
		db: db,
	}
}

func (s SiteJobStorage) IsAlreadyParsed(ctx context.Context, parseSiteJobID string) (bool, error) {
	var count int
	err := s.db.QueryStruct(ctx, &count, `
SELECT COUNT(*) FROM web_parse_page WHERE payload->>'siteJobID' = $1 AND processed_at IS NOT NULL;
`, parseSiteJobID)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
