-- +goose Up
-- +goose StatementBegin
alter table domain.scenario
	add column title text not null default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table domain.scenario
    drop column title;
-- +goose StatementEnd
