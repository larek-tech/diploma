package sitejob

import "context"

type Storage struct {
	db db
}

func New(db db) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) IsAlreadyParsed(ctx context.Context, parseSiteJobID string) (bool, error) {
	var count int
	//SELECT COUNT(*) FROM web_parse_page WHERE payload->>'payload'->>'SiteID' = '546849b4-7052-46af-8043-b9852b074afd' AND processed_at IS NOT NULL;
	err := s.db.QueryStruct(ctx, &count, `
select
	count(id)
from web_parse_page
where
	payload -> 'payload' ->> 'ID' = $1
limit 1;
`, parseSiteJobID)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (s Storage) GetProcessedPageCount(ctx context.Context, parseSiteJobID string) (int, error) {
	var count int
	err := s.db.QueryStruct(ctx, &count, `
SELECT
    COUNT(id) AS id_count
FROM
    web_parse_page
WHERE
    payload -> 'metadata' ->> 'siteJobID' = $1 AND
	processed_at IS NOT NULL;
`, parseSiteJobID)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return count, nil
	}
	return 0, nil
}

func (s Storage) GetUnprocessedPageCount(ctx context.Context, parseSiteJobID string) (int, error) {
	var count int
	err := s.db.QueryStruct(ctx, &count, `
SELECT
    COUNT(id) AS id_count
FROM
    web_parse_page
WHERE
    payload -> 'metadata' ->> 'siteJobID' = $1 AND
	processed_at IS NULL;
`, parseSiteJobID)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return count, nil
	}
	return 0, nil
}
