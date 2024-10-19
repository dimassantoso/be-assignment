-- +goose Up
-- +goose StatementBegin
INSERT INTO "payment_methods" ("name") VALUES ('Cash'), ('Credit Card'), ('Bank Transfer');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "payment_methods" RESTART IDENTITY;
-- +goose StatementEnd
