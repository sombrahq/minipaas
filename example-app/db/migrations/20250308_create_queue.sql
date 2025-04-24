-- migrate:up

SELECT minipaas_queue_create('example_queue');


-- migrate:down
SELECT minipaas_queue_delete('example_queue');
