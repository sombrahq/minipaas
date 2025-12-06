---
title: Workflows
summary: Practical end-to-end workflows using the MiniPaaS CLI.
---

## Overview

This page shows real-world workflows that combine multiple MiniPaaS CLI commands.  
These are the patterns most teams use for:

- local development
- staging / production rollouts
- CI pipelines
- routing updates
- job/worker/cron management
- secret/config provisioning

Each workflow assumes you operate inside (or pass `--env`) an environment directory containing `minipaas.yaml`.

---

# 1. Initializing an Environment

Create a new environment and register its Compose files:

```bash
minipaas code init --env dev \
  --files compose.yaml compose.build.yaml
````

This generates or updates:

* `dev/minipaas.yaml`
* environment-aware compose wiring

You can run this again whenever you change your Compose layout.

---

# 2. Defining Service Behaviors

MiniPaaS encourages **explicit service semantics** early in the development cycle.

## Add routing for HTTP services

```bash
minipaas code route --env dev \
  http://localhost:8000 \
  api:8080
```

This adds the necessary Caddy labels to the `api` service.

## Add workers

```bash
minipaas code worker --env dev emails_worker notifications_worker
```

## Add a one-off migration job

```bash
minipaas code job --env dev migrate
```

## Add a cron-triggered cleanup task

```bash
minipaas code cron --env dev cleanup
```

All these commands patch the correct Compose file in place.

---

# 3. Managing Secrets & Configs

Define secrets/configs and attach them to specific services.

## Secrets

```bash
echo postgres | minipaas secret create \
  --env dev \
  --name postgres_password \
  --for postgres \
  --for api \
  --for migrate \
  --for emails_worker
```

## Configs

```bash
minipaas config create \
  --env dev \
  --name app.json \
  ./configs/app.json \
  --for api --for emails_worker
```

Secrets/configs are created in Swarm and Compose references are patched automatically.

---

# 4. Building Images

Build all images referenced in the environment:

```bash
minipaas deploy build --env dev
```

With more output:

```bash
minipaas deploy build --verbose --env dev
```

CI pipelines often run:

```bash
minipaas deploy build --env dev --verbose
```

(Then push, depending on your Compose registry config.)

---

# 5. Deploying & Rolling Out

Apply updates with a controlled rollout:

```bash
minipaas deploy rollout --env dev
```

This performs:

* image tag resolution
* service-level updates
* rollout sequencing

Verbose mode:

```bash
minipaas deploy rollout --verbose --env dev
```

---

# 6. Updating Routing (Caddy)

Whenever routing-related scaffolding changes, reapply routing:

```bash
minipaas deploy routing --env dev
```

This regenerates and updates routing config inside the stack.

---

# 7. Running Jobs Manually

Job services created with `code job` are typically run by CI or deploy workflows.

Example: run a migration job before rollout:

```bash
minipaas deploy build --env dev
minipaas deploy rollout --env dev
# then run job via swarm invocation or CI automation (varies by stack)
```

MiniPaaS CLI itself does not “run jobs” directly — Swarm orchestrates them after deployment.

---

# 8. Cron Execution

Services marked via `code cron` run automatically through `swarm-cronjob`.

Make sure the Ansible role is used to provision the cluster so that:

* `swarm-cronjob` is installed
* cron schedules defined in Compose take effect

No additional CLI step is needed after scaffolding.

---

# 9. CI / CD Pipeline Pattern

A typical CI pipeline:

```bash
# Install dependencies / CLI
minipaas code init --env prod --files compose.yaml compose.build.yaml

# Provision secrets/configs (common in CI)
echo "${POSTGRES_PASSWORD}" | minipaas secret create \
  --env prod --name postgres_password --for postgres --for api

# Build & rollout
minipaas deploy build --env prod --verbose
minipaas deploy rollout --env prod --verbose

# Apply updated routing
minipaas deploy routing --env prod
```

This keeps CI pipelines predictable and environment definitions explicit.

---

# Summary

MiniPaaS workflows combine:

* `code` → service semantics
* `secret` / `config` → payload & metadata management
* `deploy build` → build images
* `deploy rollout` → update services
* `deploy routing` → refresh ingress config

Use this page as a pattern reference; for command-level details see the **Commands** page.
