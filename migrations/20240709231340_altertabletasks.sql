-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
    ADD COLUMN started_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    ADD COLUMN completed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    ADD COLUMN total_duration INTERVAL DEFAULT '0 seconds';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
    DROP COLUMN started_at,
    DROP COLUMN completed_at,
    DROP COLUMN total_duration;
-- +goose StatementEnd
