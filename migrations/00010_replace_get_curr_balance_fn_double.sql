-- +goose Up
-- +goose StatementBegin
DROP FUNCTION get_current_user_balance;
CREATE OR REPLACE FUNCTION get_current_user_balance(p_user_id INT)
RETURNS DOUBLE PRECISION AS $$
DECLARE
    v_last_checkpoint_balance DOUBLE PRECISION;
    v_recent_sum DOUBLE PRECISION;
BEGIN
    -- Get the last checkpoint balance
    SELECT running_balance INTO v_last_checkpoint_balance
    FROM transaction_table
    WHERE user_id = p_user_id
    AND is_balance_checkpoint = TRUE
    ORDER BY processed_at DESC
    LIMIT 1;
    
    -- Sum transactions after the last checkpoint
    SELECT COALESCE(SUM(amount), 0) INTO v_recent_sum
    FROM transaction_table
    WHERE user_id = p_user_id
    AND (
        is_balance_checkpoint = FALSE OR
        processed_at > (
            SELECT MAX(processed_at)
            FROM transaction_table
            WHERE user_id = p_user_id
            AND is_balance_checkpoint = TRUE
        )
    );
    
    RETURN COALESCE(v_last_checkpoint_balance, 0) + v_recent_sum;
END;
$$ LANGUAGE plpgsql
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION get_current_user_balance;
-- +goose StatementEnd
