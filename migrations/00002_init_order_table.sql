-- +goose Up
-- +goose StatementBegin
CREATE TYPE eval_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TABLE IF NOT EXISTS order_table (
  order_id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  order_num VARCHAR(32) UNIQUE NOT NULL,
  amount INT NOT NULL,
  points_eval_status EVAL_STATUS NOT NULL DEFAULT 'NEW',
  uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES user_table(user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_table;
-- +goose StatementEnd
