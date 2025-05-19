-- +goose Up
-- +goose StatementBegin
alter table domain.scenario
    add column context_size int not null default 8192;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table domain.scenario
    drop column context_size;
-- +goose StatementEnd
