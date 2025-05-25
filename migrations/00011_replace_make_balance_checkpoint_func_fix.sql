-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION make_balance_checkpoint(
    p_user_id INT, 
    p_nth INT
) RETURNS VOID AS $$
DECLARE
    v_last_checkpoint_id BIGINT;
    v_latest_id BIGINT;
    v_last_balance DOUBLE PRECISION;
    v_transaction_count INT;
    v_new_balance DOUBLE PRECISION;
    v_max_batch_id INT;
BEGIN
    -- Find the last checkpoint
    SELECT transaction_id, running_balance INTO v_last_checkpoint_id, v_last_balance
    FROM transaction_table
    WHERE user_id = p_user_id 
    AND is_balance_checkpoint = TRUE
    ORDER BY processed_at DESC
    LIMIT 1;
    
    -- Count transactions since last checkpoint
    SELECT COUNT(*) INTO v_transaction_count
    FROM transaction_table
    WHERE user_id = p_user_id 
    AND (v_last_checkpoint_id IS NULL OR transaction_id > v_last_checkpoint_id);
    
    -- If we've reached N transactions, calculate new balance
    IF v_transaction_count >= p_nth THEN
        -- Calculate sum since last checkpoint
        SELECT COALESCE(SUM(amount), 0) INTO v_new_balance
        FROM transaction_table
        WHERE user_id = p_user_id 
        AND (v_last_checkpoint_id IS NULL OR transaction_id > v_last_checkpoint_id);
        
        -- Add to previous balance (or start from 0)
        v_new_balance := COALESCE(v_last_balance, 0) + v_new_balance;

        -- Get the next batch ID
        SELECT COALESCE(MAX(batch_id), 0) + 1 INTO v_max_batch_id
        FROM transaction_table 
        WHERE user_id = p_user_id;

        -- Find the latest transaction ID in this batch
        SELECT MAX(transaction_id) INTO v_latest_id
        FROM transaction_table
        WHERE user_id = p_user_id 
        AND (v_last_checkpoint_id IS NULL OR transaction_id > v_last_checkpoint_id);
            
        UPDATE transaction_table
        SET 
            batch_id = v_max_batch_id,
            running_balance = CASE 
                WHEN transaction_id = v_latest_id THEN v_new_balance
                ELSE NULL
            END,
            is_balance_checkpoint = (transaction_id = v_latest_id)
        WHERE user_id = p_user_id 
        AND (v_last_checkpoint_id IS NULL OR transaction_id > v_last_checkpoint_id);
        
    END IF;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION make_balance_checkpoint;
-- +goose StatementEnd
