-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "payment_methods" (
     "id" SERIAL PRIMARY KEY,
     "name" VARCHAR(50) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_methods;
-- +goose StatementEnd
