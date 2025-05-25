-- +goose Up
-- +goose StatementBegin
DROP FUNCTION get_current_user_balance;
CREATE OR REPLACE FUNCTION get_current_user_balance(p_user_id INT)
RETURNS DOUBLE PRECISION AS $$
DECLARE
    v_last_checkpoint_balance DOUBLE PRECISION;
    v_recent_sum DOUBLE PRECISION;
    v_last_checkpoint_ts TIMESTAMPTZ;
BEGIN
    -- Get the last checkpoint balance
    SELECT running_balance, processed_at  
    INTO v_last_checkpoint_balance, v_last_checkpoint_ts
    FROM transaction_table
    WHERE user_id = p_user_id
    AND is_balance_checkpoint = TRUE
    ORDER BY processed_at DESC
    LIMIT 1;

    -- Sum ONLY transactions after the last checkpoint (or all if no checkpoint exists)
    IF v_last_checkpoint_ts IS NOT NULL THEN
      SELECT COALESCE(SUM(amount), 0) INTO v_recent_sum
      FROM transaction_table
      WHERE user_id = p_user_id
      AND processed_at > v_last_checkpoint_ts
      AND is_balance_checkpoint = FALSE;
    ELSE
      SELECT COALESCE(SUM(amount), 0) INTO v_recent_sum
      FROM transaction_table
      WHERE user_id = p_user_id;
    END IF;
    
    RETURN COALESCE(v_last_checkpoint_balance, 0) + v_recent_sum;
END;
$$ LANGUAGE plpgsql
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION get_current_user_balance;
-- +goose StatementEnd
