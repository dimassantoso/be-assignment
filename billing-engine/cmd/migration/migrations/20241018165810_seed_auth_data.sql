-- +goose Up
-- +goose StatementBegin
-- initial auth email: admin@example.com & password: P@SSword1234
INSERT INTO "auths" (email, password) VALUES ('admin@example.com', '$2a$10$kl5A0MKtKdjCcj.gPyqWwOFIs96oWhj5Ryqikr9AWqXrN1wfmk3eS');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "auths" RESTART IDENTITY;
-- +goose StatementEnd
