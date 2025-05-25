-- +goose Up
-- +goose StatementBegin
CREATE INDEX u_batch_idx
ON transaction_table(user_id, batch_id);

CREATE INDEX u_time_idx
ON transaction_table(user_id, processed_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX u_batch_idx;
DROP INDEX u_time_idx;
-- +goose StatementEnd
