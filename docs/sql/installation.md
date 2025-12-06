---
title: Installation
summary: How to install the MiniPaaS SQL functions and tables into a PostgreSQL database.
---

## Overview

The `minipaas-sql/` package contains SQL files that define durable queue and stream primitives for PostgreSQL.  
These scripts are **dbmate-ready**, making them easy to integrate into migration workflows used in production environments.

The SQL schema can be installed using:

- `psql`
- dbmate
- a project’s migration system (Flyway, Goose, Sqitch, Django migrations, Prisma, Liquibase)
- CI/CD pipelines running SQL steps

---

## Requirements

To install the SQL component, you need:

- A PostgreSQL instance  
- A database user with permission to create tables and functions  
- Either `psql`, dbmate, or your migration system  

---

## Installation With psql

```bash
psql "$DATABASE_URL" -f minipaas-sql/tables.sql
psql "$DATABASE_URL" -f minipaas-sql/queues.sql
psql "$DATABASE_URL" -f minipaas-sql/streams.sql
````

Each file defines a portion of the durable queue and stream functionality.

---

## Installation With dbmate

The SQL files are structured so they can be used directly as dbmate migrations.

You may either:

### Option 1 — Copy SQL files directly into your dbmate migrations directory

Example:

```
db/migrations/
  20250203120000_create_queues.sql
  20250203120001_create_streams.sql
  20250203120002_supporting_functions.sql
```

Then run:

```bash
dbmate up
```

### Option 2 — Reference them as external migration steps

(Some teams prefer symlinks or include directives depending on their CI/migration layout.)

Both approaches allow dbmate to track schema versions, enforce ordering, and support rollback (if you add corresponding down-migrations).

---

## Installation in CI or Deployment Pipelines

Common patterns include:

* Executing dbmate migrations during deployment
* Running `psql -f` commands before starting workers
* Including SQL in Docker image entrypoints
* Using a GitOps-style migration controller

As long as the SQL files are applied before queues and streams are used, applications will be able to enqueue, dequeue, publish, and consume.

---

## Using Other Migration Systems

The SQL can be imported directly into:

* Flyway
* Goose
* Sqitch
* Django migrations (as raw SQL)
* Prisma SQL migrations
* Liquibase XML/SQL migrations

This allows version control, rollbacks, and integration with existing schema evolution workflows.

---

## Verifying Installation

You can verify the installation using any SQL client by checking for:

* queue tables
* stream tables
* functions supporting enqueue/dequeue
* functions supporting publish/consume

All objects are created in ordinary PostgreSQL schemas.

---

## Summary

To install the SQL layer:

1. Apply the SQL files in `minipaas-sql/` using `psql`, dbmate, or your migration tool.
2. Optionally integrate the SQL scripts into your project’s migrations.
3. Run migrations before workers or services start using queues or streams.
4. Begin enqueueing tasks and publishing events from your applications.

Because the SQL implementation is portable and dbmate-compatible, it fits easily into existing operational workflows.

