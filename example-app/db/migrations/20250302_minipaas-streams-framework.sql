-- migrate:up

-- =============================================
-- Create a new stream table for a given stream name
-- =============================================
CREATE OR REPLACE FUNCTION minipaas_stream_create(p_stream_name text)
    RETURNS void AS
$$
DECLARE
    table_name text := p_stream_name;
BEGIN
    IF p_stream_name !~ '^[a-z_][a-z0-9_]*$' THEN
        RAISE EXCEPTION 'Invalid stream table name: %', p_stream_name;
    END IF;

    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I (
            id BIGSERIAL PRIMARY KEY,
            payload JSONB NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT now()
        )', table_name);
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- Delete a stream table by name
-- =============================================
CREATE OR REPLACE FUNCTION minipaas_stream_delete(p_stream_name text)
    RETURNS void AS
$$
DECLARE
    table_name text := p_stream_name;
BEGIN
    IF p_stream_name !~ '^[a-z_][a-z0-9_]*$' THEN
        RAISE EXCEPTION 'Invalid stream table name: %', p_stream_name;
    END IF;

    EXECUTE format('DROP TABLE IF EXISTS %I', table_name);
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- Publish a message to a specified stream.
-- Returns the inserted message.
-- =============================================
CREATE OR REPLACE FUNCTION minipaas_stream_publish(p_stream_name text, p_payload jsonb)
    RETURNS TABLE(
                     id bigint,
                     payload jsonb,
                     created_at timestamptz
                 )
AS
$$
DECLARE
    table_name text := p_stream_name;
    new_id bigint;
    result_payload jsonb;
    result_created_at timestamptz;
    notify_channel text := p_stream_name;  -- Use the stream name directly as channel.
BEGIN
    IF p_stream_name !~ '^[a-z_][a-z0-9_]*$' THEN
        RAISE EXCEPTION 'Invalid stream table name: %', p_stream_name;
    END IF;

    EXECUTE format('
        INSERT INTO %I (payload)
        VALUES ($1)
        RETURNING id, payload, created_at
    ', table_name)
        INTO new_id, result_payload, result_created_at
        USING p_payload;

    PERFORM pg_notify(notify_channel, json_build_object('id', new_id, 'payload', result_payload)::text);

    id := new_id;
    payload := result_payload;
    created_at := result_created_at;
    RETURN NEXT;
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- Subscribe to a stream: return N messages with id greater than p_last_id.
-- =============================================
CREATE OR REPLACE FUNCTION minipaas_stream_consume(p_stream_name text, p_last_id bigint, p_batch_size int)
    RETURNS TABLE(
                     id bigint,
                     payload jsonb,
                     created_at timestamptz
                 )
AS
$$
DECLARE
    table_name text := p_stream_name;
BEGIN
    IF p_stream_name !~ '^[a-z_][a-z0-9_]*$' THEN
        RAISE EXCEPTION 'Invalid stream table name: %', p_stream_name;
    END IF;

    IF p_batch_size <= 0 THEN
        RAISE EXCEPTION 'Batch size must be greater than zero';
    END IF;

    RETURN QUERY EXECUTE format('
        SELECT id, payload, created_at
        FROM %I
        WHERE id > $1
        ORDER BY id ASC
        LIMIT $2
    ', table_name)
        USING p_last_id, p_batch_size;
END;
$$ LANGUAGE plpgsql;

-- migrate:down
