---
title: Scaffolding
summary: Add routing, workers, jobs, and cron behavior by modifying your Compose files.
---

## Overview

The `minipaas code` commands **modify your Compose files in place** to attach operational behaviors such as routing, background workers, one-off jobs, and cron tasks.  
Changes are made only to the Compose files declared in `minipaas.yaml` under the selected `--env` directory.

Scaffolding is **idempotent** — re-running commands updates what’s needed and never duplicates or corrupts Compose definitions.

---

## How Scaffolding Works

1. The CLI loads all Compose files listed under `compose:` in `minipaas.yaml`.
2. It finds the *first* file where each service is defined.
3. It applies the required modifications (labels, modes, restart settings, annotations).
4. It writes the updated Compose file back to disk.

Only the relevant Compose file for each service is modified — nothing else is touched.

---

## Initialize an Environment

Create or update an environment directory, generate a `minipaas.yaml`, and wire Compose files:

```bash
minipaas code init --env dev \
  --files compose.yaml compose.build.yaml
````

This sets up a ready-to-use MiniPaaS environment.

---

## Route a Service (Caddy Ingress)

Use `code route` to expose a service through Caddy.

```bash
minipaas code route --env dev \
  http://localhost:8000 \
  api:8080
```

Arguments:

1. **Public URL** (scheme + host + optional port)
2. **Service mapping** in the form `service:internal-port`

This adds Caddy labels to the service so that HTTP traffic is routed correctly inside the Swarm stack.

---

## Background Workers

Convert one or more services into **long-running background workers**:

```bash
minipaas code worker --env dev worker emails_worker
```

Worker scaffolding typically:

* Removes exposed ports
* Ensures deterministic restart behavior
* Marks the service as a worker for downstream tooling (e.g., logs, metrics, jobs)
* Keeps everything else in your Compose file intact

Use workers for queue consumers, stream processors, or any always-running backend task.

---

## One-Off Jobs

Mark a service as a **run-once job**:

```bash
minipaas code job --env dev migrate
```

Job services are:

* Expected to **run once and exit**
* Typically executed during deployments or CI pipelines
* Not scaled or routed like normal services

This is ideal for database migrations, cleanup tasks, and one-time scripts.

---

## Cron Tasks (via swarm-cronjob)

Mark a service as a **scheduled cron task**:

```bash
minipaas code cron --env dev cleanup
```

The schedule itself is defined in your Compose or environment configuration; the CLI simply attaches the required labels so that `swarm-cronjob` can trigger it.

The MiniPaaS Ansible role installs and configures `swarm-cronjob` on the cluster, so cron tasks execute automatically.

---

## Best Practices

* Put scaffolding under version control so compose changes are visible in PRs.
* Scaffold early — routing, workers, cron, and jobs are **service semantics**, not deployment-time details.
* Running scaffolding repeatedly is safe; commands are idempotent.
* If the CLI reports “service not found,” check whether the service appears in one of the compose files listed in `minipaas.yaml`.

