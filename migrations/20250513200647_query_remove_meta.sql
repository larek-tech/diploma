-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
alter table chat.query
    drop column metadata;
alter table chat.query
    drop column source_ids;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table chat.query
    add column metadata jsonb not null default '{}'::jsonb;
alter table chat.query
    add column source_ids bigint[] not null default array[]::bigint[];
-- +goose StatementEnd
