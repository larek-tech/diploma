-- +goose Up
-- +goose StatementBegin
alter table domain.domain
    add column scenario_ids bigint[] not null default array[]::bigint[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table domain.domain
    drop column scenario_ids;
-- +goose StatementEnd
