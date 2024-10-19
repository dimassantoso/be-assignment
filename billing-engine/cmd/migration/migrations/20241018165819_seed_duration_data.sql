-- +goose Up
-- +goose StatementBegin
INSERT INTO "durations" ("week", "interest")
VALUES (50, 10), (36, 15);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "durations" RESTART IDENTITY;
-- +goose StatementEnd
