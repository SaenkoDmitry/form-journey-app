-- +goose Up
-- +goose StatementBegin
UPDATE exercise_types SET units = 'reps,weight' WHERE id in (1,2,3,5,6,7,8,9,10,11,12,13,14,15,16,17);
UPDATE exercise_types SET units = 'minutes' WHERE id in (4,18,19,20,21);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE exercise_types SET units = null WHERE id in (1,2,3,5,6,7,8,9,10,11,12,13,14,15,16,17);
UPDATE exercise_types SET units = null WHERE id in (4,18,19,20,21);
-- +goose StatementEnd
