-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ALTER COLUMN passport TYPE VARCHAR(11),
    ALTER COLUMN passport SET NOT NULL,
    ADD CONSTRAINT passport_unique UNIQUE (passport),
    ALTER COLUMN surname DROP NOT NULL,
    ALTER COLUMN name DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    ALTER COLUMN passport TYPE VARCHAR(255),
    DROP CONSTRAINT IF EXISTS passport_unique,
    ALTER COLUMN surname SET NOT NULL,
    ALTER COLUMN name SET NOT NULL;
-- +goose StatementEnd
