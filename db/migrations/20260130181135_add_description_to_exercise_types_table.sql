-- +goose Up
-- +goose StatementBegin
ALTER TABLE exercise_types ADD COLUMN description TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE exercise_types DROP COLUMN description;
-- +goose StatementEnd
