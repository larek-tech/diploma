-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- Enable vector extension for pgvector
CREATE EXTENSION IF NOT EXISTS vector;

-- Create main tables
CREATE TABLE IF NOT EXISTS sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    type SMALLINT NOT NULL,
    credentials BYTEA
);

CREATE TABLE IF NOT EXISTS sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID REFERENCES sources(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    metadata JSONB,
    raw TEXT,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID REFERENCES sources(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Vector table for embeddings
CREATE TABLE IF NOT EXISTS chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    index INT NOT NULL,
    document_id UUID REFERENCES documents(id) ON DELETE CASCADE,
    content TEXT,
    embeddings VECTOR(1024)
);

-- Index for vector similarity search
CREATE INDEX ON chunks USING ivfflat (embeddings vector_cosine_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS chunks;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS pages;
DROP TABLE IF EXISTS sites;
DROP TABLE IF EXISTS sources;
DROP EXTENSION IF EXISTS vector;
-- +goose StatementEnd