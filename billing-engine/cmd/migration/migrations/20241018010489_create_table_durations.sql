-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "durations" (
   "id" SERIAL PRIMARY KEY,
   "week" INT NOT NULL,
   "interest" DECIMAL(5, 2) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS durations;
-- +goose StatementEnd
