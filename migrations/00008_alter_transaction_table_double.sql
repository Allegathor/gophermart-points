-- +goose Up
-- +goose StatementBegin
ALTER TABLE transaction_table
ALTER COLUMN amount TYPE DOUBLE PRECISION,
ALTER COLUMN running_balance TYPE DOUBLE PRECISION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transaction_table
ALTER COLUMN amount TYPE INT,
ALTER COLUMN running_balance TYPE INT;
-- +goose StatementEnd
