-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "loans" (
    "id" SERIAL PRIMARY KEY,
    "borrower_id" INT REFERENCES "borrowers"("id") ON DELETE CASCADE,
    "duration_id" INT REFERENCES "durations"("id") ON DELETE SET NULL,
    "principal_amount" DECIMAL(12, 2) NOT NULL,
    "interest_amount" DECIMAL(12, 2) NOT NULL,
    "total_amount" DECIMAL(12, 2) NOT NULL,
    "last_payment_date" TIMESTAMPTZ(6),
    "outstanding" DECIMAL(12, 2),
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS loans;
-- +goose StatementEnd
