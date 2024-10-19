-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS billings (
    "id" SERIAL PRIMARY KEY,
    "loan_id" INT REFERENCES "loans"("id") ON DELETE CASCADE,
    "week" INT NOT NULL,
    "due_date" DATE NOT NULL,
    "amount_due" DECIMAL(12, 2) NOT NULL,
    "payment_date" TIMESTAMPTZ(6),
    "payment_method_id" INT REFERENCES "payment_methods"("id"),
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS billings;
-- +goose StatementEnd


