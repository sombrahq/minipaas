-- migrate:up
CREATE TABLE records
(
    id   SERIAL PRIMARY KEY,
    data TEXT NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS records;
