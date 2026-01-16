-- +goose Up
-- +goose StatementBegin
ALTER TABLE exercise_types ADD COLUMN units TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE exercise_types DROP COLUMN units;
-- +goose StatementEnd
