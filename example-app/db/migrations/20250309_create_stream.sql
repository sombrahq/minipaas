-- migrate:up

SELECT minipaas_stream_create('example_stream');

CREATE TABLE consumer
(
    id            uuid PRIMARY KEY,
    last_event_id BIGINT,
    updated_at    TIMESTAMPTZ DEFAULT now()
);

-- migrate:down

DROP TABLE IF EXISTS consumers;
SELECT minipaas_stream_delete('example_stream');
