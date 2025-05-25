-- +goose Up
-- +goose StatementBegin
ALTER TABLE order_table
ALTER COLUMN amount TYPE DOUBLE PRECISION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE order_table
ALTER COLUMN amount TYPE INT;
-- +goose StatementEnd
