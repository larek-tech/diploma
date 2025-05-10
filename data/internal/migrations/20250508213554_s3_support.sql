-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS object_storage (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id uuid REFERENCES sources(id) ON DELETE CASCADE,
    config JSON
);

CREATE TABLE IF NOT EXISTS object_storage_files (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    object_storage_id uuid REFERENCES object_storage(id) ON DELETE CASCADE,
    raw_object_id text NOT NULL, -- reference to object inside service wide object storage
    content_type text NOT NULL, -- content type of the object
    content_size BIGINT NOT NULL, -- size of the object
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS object_storage_files;
DROP TABLE IF EXISTS object_storage;

-- +goose StatementEnd
