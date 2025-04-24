-- migrate:up

-- =============================================
-- Create a queue table with indexing and notifications
-- =============================================
CREATE OR REPLACE FUNCTION minipaas_queue_create(p_queue_name text)
    RETURNS void AS
$$
DECLARE
    table_name text := p_queue_name;
BEGIN
    -- Validate the table name to prevent SQL injection and ensure correct formatting.
    IF p_queue_name !~ '^[a-z_][a-z0-9_]*$' THEN
        RAISE EXCEPTION 'Invalid queue table name: %', p_queue_name;
    END IF;

    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I (
            id SERIAL PRIMARY KEY,
            payload JSONB NOT NULL,
            status VARCHAR(20) NOT NULL DEFAULT ''pending'',
            created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
            processed_at TIMESTAMPTZ
        )', table_name);

    EXECUTE format('
        CREATE INDEX IF NOT EXISTS idx_%I_status_created_at
        ON %I (status, created_at)', table_name, table_name);
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- Enqueue a new job and notify listeners
-- =============================================
CREATE OR REPLACE FUNCTION minipaas_queue_enqueue(p_queue_name text, p_payload jsonb)
    RETURNS TABLE(
                     id         int,
                     payload    jsonb,
                     status     varchar,
                     created_at timestamptz
                 )
AS
$$
DECLARE
    table_name text := p_queue_name;
    v_new_id       int;
    v_payload      jsonb;
    v_status       varchar;
    v_created_at   timestamptz;
    notify_channel text := p_queue_name;  -- using the table name as the notification channel
BEGIN
    IF p_queue_name !~ '^[a-z_][a-z0-9_]*$' THEN
        RAISE EXCEPTION 'Invalid queue table name: %', p_queue_name;
    END IF;

    EXECUTE format('
        INSERT INTO %I (payload)
        VALUES ($1)
        RETURNING id, payload, status, created_at
    ', table_name)
        INTO v_new_id, v_payload, v_status, v_created_at
        USING p_payload;

    PERFORM pg_notify(notify_channel, v_new_id::text);

    id := v_new_id;
    payload := v_payload;
    status := v_status;
    created_at := v_created_at;
    RETURN NEXT;
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- Dequeue multiple jobs (batch processing)
-- =============================================
CREATE OR REPLACE FUNCTION minipaas_queue_dequeue(p_queue_name text, p_batch_size int)
    RETURNS TABLE(
                     id         int,
                     payload    jsonb,
                     status     varchar,
                     created_at timestamptz
                 )
AS
$$
DECLARE
    table_name text := p_queue_name;
BEGIN
    IF p_queue_name !~ '^[a-z_][a-z0-9_]*$' THEN
        RAISE EXCEPTION 'Invalid queue table name: %', p_queue_name;
    END IF;

    IF p_batch_size <= 0 THEN
        RAISE EXCEPTION 'Batch size must be greater than zero';
    END IF;

    RETURN QUERY EXECUTE format('
        WITH jobs AS (
            SELECT id FROM %I
            WHERE status = ''pending''
            ORDER BY created_at
            FOR UPDATE SKIP LOCKED
            LIMIT $1
        )
        UPDATE %I
        SET status = ''processing''
        WHERE id IN (SELECT id FROM jobs)
        RETURNING id, payload, status, created_at
    ', table_name, table_name)
        USING p_batch_size;
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- Acknowledge multiple messages
-- =============================================
CREATE OR REPLACE FUNCTION minipaas_queue_ack(p_queue_name text, p_job_ids int[])
    RETURNS void AS
$$
DECLARE
    table_name text := p_queue_name;
BEGIN
    IF p_queue_name !~ '^[a-z_][a-z0-9_]*$' THEN
        RAISE EXCEPTION 'Invalid queue table name: %', p_queue_name;
    END IF;

    EXECUTE format('
        UPDATE %I
        SET status = ''done'', processed_at = now()
        WHERE id = ANY($1)
    ', table_name)
        USING p_job_ids;
END;
$$ LANGUAGE plpgsql;

-- migrate:down
