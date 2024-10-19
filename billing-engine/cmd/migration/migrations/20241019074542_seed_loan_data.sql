-- +goose Up
-- +goose StatementBegin
INSERT INTO loans (borrower_id, duration_id, principal_amount, interest_amount, total_amount, last_payment_date, outstanding) VALUES (1, 1, 5000000.00, 500000.00, 5500000.00, '2024-09-21 00:00:00.000000 +00:00', 1650000.00);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE billing RESTART IDENTITY;
TRUNCATE TABLE loans RESTART IDENTITY CASCADE;
-- +goose StatementEnd
