-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE chunks ALTER COLUMN embeddings SET TYPE vector(8192);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE chunks ALTER COLUMN embeddings SET TYPE vector(1024);
-- +goose StatementEnd
