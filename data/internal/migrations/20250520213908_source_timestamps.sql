-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE sources ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE sources ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Create a trigger function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create a trigger on the sources table
CREATE TRIGGER update_sources_updated_at
BEFORE UPDATE ON sources
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TRIGGER IF EXISTS update_sources_updated_at ON sources;
DROP FUNCTION IF EXISTS update_updated_at_column();
ALTER TABLE sources DROP COLUMN created_at;
ALTER TABLE sources DROP COLUMN updated_at;
-- +goose StatementEnd