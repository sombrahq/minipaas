---
title: Usage
summary: How applications enqueue, dequeue, publish, and consume using MiniPaaS SQL.
---

## Overview

The `minipaas-sql/` package provides SQL primitives for two core behaviors:

- **Durable queues** — for background workers and asynchronous jobs  
- **Durable streams** — for event-style processing with ordered offsets  

Applications interact with these features directly through SQL.  
This page shows the typical usage patterns for producers and consumers.

---

## Durable Queues

Queues allow applications to store units of work and process them through workers.

### Enqueueing Work

Workers or services add tasks to a queue using the enqueue function:

```sql
SELECT queue_enqueue('my_queue', '{"task": "process-user", "id": 42}');
````

Payloads are stored as JSON, making them flexible for application needs.

---

### Dequeuing Tasks

A worker retrieves the next available task:

```sql
SELECT * FROM queue_dequeue('my_queue');
```

This returns at most one task for the worker to process.

Workers typically run this in a loop with backoff or timers, depending on service behavior.

---

### Acknowledging Completion

When the worker completes the task:

```sql
SELECT queue_ack(<task_id>);
```

This removes the task from the queue.

---

## Durable Streams

Streams allow producers to publish ordered events and consumers to track offsets.

### Publishing an Event

```sql
SELECT stream_publish('activity_log', '{"event": "signup", "user": 12}');
```

Each event receives a strictly increasing offset.

---

### Reading Events

Consumers fetch new events by specifying their last processed offset:

```sql
SELECT * FROM stream_consume('activity_log', <last_offset>);
```

The consumer updates its offset in application storage once processing is complete.

---

## Worker & Consumer Patterns

### Worker Processes (Queues)

Typical workflow:

1. `queue_dequeue` to claim a task
2. Process the payload
3. `queue_ack` to acknowledge
4. Repeat

This model fits background workers, asynchronous jobs, and retry logic defined at the application level.

---

### Stream Consumers

Typical workflow:

1. Track a consumer offset
2. `stream_consume` to fetch new events
3. Process events in order
4. Update the offset
5. Repeat

This suits audit trails, notifications, analytics pipelines, and cross-service event delivery.

---

## Using Queues and Streams With MiniPaaS CLI

Applications deployed with MiniPaaS often use:

* **workers** → consuming from queues
* **jobs** → processing a queue once and exiting
* **cron services** → publishing to queues or streams on schedule

These service types are scaffolded by the CLI and interact directly with PostgreSQL using your application code.

The SQL layer remains unchanged across environments; only your application configuration varies.

---

## Summary

Applications use queues and streams by calling a small set of SQL functions:

* enqueue and dequeue tasks
* acknowledge processing
* publish and consume ordered events

These primitives integrate naturally with worker, job, and cron services deployed via MiniPaaS and can support a wide range of application patterns without additional infrastructure.

