-- +goose Up
-- +goose StatementBegin
INSERT INTO borrowers (name, email) VALUES ('John Doe', 'john.doe@gmail.com'), ('Michael Doe', 'michael.doe@gmail.com');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE borrowers RESTART IDENTITY;
-- +goose StatementEnd
