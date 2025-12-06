---
title: SQL Overview
summary: PostgreSQL functions for durable queues and streams used by MiniPaaS applications.
---

## Overview

The `minipaas-sql/` component provides a set of PostgreSQL tables and functions that implement durable queues and durable streams.  
These capabilities support background workers, scheduled tasks, event-driven flows, and small-scale message distribution patterns.

The design focuses on simplicity, portability, and clarity.  
All logic is implemented using standard PostgreSQL features and can be deployed in any PostgreSQL database.

---

## Features

### Durable Queues
Provides functions and tables for applications that need to enqueue units of work and process them through workers.  
Concurrency, safe dequeueing, and acknowledgement behavior are handled through PostgreSQL’s locking and transactional guarantees.

### Durable Streams
Allows applications to publish ordered events and process them through consumer offsets.  
Streams provide a lightweight building block for audit trails, activity logs, or inter-service event consumption.

### Pure SQL Implementation
All logic is defined in SQL files.  
There are no extensions, background daemons, or external dependencies.  
This keeps the system easy to adopt, maintain, inspect, and evolve.

### Flexible Deployment
The SQL can be installed into:

- the Postgres instance provided by your infrastructure  
- a dedicated application database  
- an existing cluster shared by multiple services  

Applications decide how to consume queues and streams through their own workers and services.

---

## Behavior

The SQL package defines:

- tables to store queued tasks and stream events  
- functions to enqueue, dequeue, publish, and read events  
- concurrency semantics using row locking  
- simple, deterministic ordering for streams  

The model is designed to integrate naturally with the MiniPaaS CLI’s worker, job, and cron workflows, while remaining entirely optional and reusable outside MiniPaaS.

---

## Start Here

- Install the SQL schema → **[Installation](installation.md)**  
- Learn how to enqueue, dequeue, publish, and consume → **[Usage](usage.md)**  
