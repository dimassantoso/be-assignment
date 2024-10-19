-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "borrowers" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(50) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "deleted_at" TIMESTAMPTZ(6) DEFAULT NULL,
    "created_at" TIMESTAMPTZ(6) DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ(6) DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS borrowers;
-- +goose StatementEnd
