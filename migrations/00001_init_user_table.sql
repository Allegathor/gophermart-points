-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_table (
  user_id SERIAL PRIMARY KEY,
  login VARCHAR(32) UNIQUE,
  passwd VARCHAR(60)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_table;
-- +goose StatementEnd
