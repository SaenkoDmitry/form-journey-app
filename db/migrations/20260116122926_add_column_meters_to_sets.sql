-- +goose Up
-- +goose StatementBegin
ALTER TABLE sets ADD COLUMN meters INT;
ALTER TABLE sets ADD COLUMN fact_meters INT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sets DROP COLUMN meters;
ALTER TABLE sets DROM COLUMN fact_meters;
-- +goose StatementEnd
