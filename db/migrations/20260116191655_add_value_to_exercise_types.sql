-- +goose Up
-- +goose StatementBegin
INSERT INTO exercise_group_types (code, name)
VALUES ('buttocks', '🍑 Ягодицы');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE exercise_group_types;
-- +goose StatementEnd
