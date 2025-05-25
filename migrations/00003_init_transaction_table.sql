-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_table (
  transaction_id BIGSERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  order_id INT NOT NULL,
  batch_id INT,
  is_balance_checkpoint BOOLEAN DEFAULT FALSE,
  running_balance INT,
  amount INT NOT NULL,
  processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES user_table(user_id),
  FOREIGN KEY (order_id) REFERENCES order_table(order_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transaction_table;
-- +goose StatementEnd

