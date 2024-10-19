-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS auths (
    "id" SERIAL PRIMARY KEY,
    "email" VARCHAR(50) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ(6) DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ(6) DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auths;
-- +goose StatementEnd
