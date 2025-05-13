-- +goose Up
-- +goose StatementBegin
alter table domain.scenario
    add column domain_id bigint not null default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table domain.scenario
    drop column domain_id;
-- +goose StatementEnd
