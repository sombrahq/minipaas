---
title: About
summary: What MiniPaaS is, why it exists, and how the pieces fit together.
---

## The Problem

Running production services requires repeatable deployments, TLS, firewalling, monitoring, secrets management, and a way to coordinate background work.  
Teams often face a difficult choice:

- **Kubernetes** — powerful but heavy for small clusters; sizable operational overhead.
- **Docker Compose** — simple but locked to a single host and lacking production workflows.
- **Raw Docker Swarm** — good primitives but few conventions or batteries included.

The gap is especially visible for small teams who want infrastructure that is:

- Understandable end-to-end
- Automatable
- CI-friendly
- Not beholden to cloud-specific control planes

## The Idea

MiniPaaS provides a small set of open components that layer onto Docker Swarm:

- **CLI** (`minipaas-cli/`): Orchestrates deployments, rollouts, routing, workers, jobs, cron, and secrets/configs — all from Compose files.
- **Ansible Role** (`minipaas-role/`): Turns plain Linux hosts into secure Swarm nodes with TLS, firewalling, monitoring, syslog, Caddy, and swarm-cronjob.
- **SQL** (`minipaas-sql/`): Adds durable queues and streams directly inside PostgreSQL to support background work and event pipelines.

Each component is standalone. You can adopt one, two, or all three.

## Why Build It?

MiniPaaS aims to reduce infrastructure friction while keeping everything transparent:

- Give small teams practical tooling that doesn’t require Kubernetes expertise.
- Provide a consistent deployment workflow based on Compose files teams already use.
- Standardize Swarm node provisioning with predictable security defaults.
- Use Postgres for queues/streams to avoid maintaining extra infrastructure.
- Keep migration paths open — nothing blocks moving to another platform later.

## Principles

- **Small, composable parts** — adopt what you need.
- **Explicit over implicit** — no hidden control planes or magic.
- **Self-hosted by default** — works offline, locally, and in CI.
- **Readable and auditable** — everything is plain Go, Ansible, SQL, and YAML.

## Learn More

- Infrastructure: [Ansible Role](role/index.md)
- Deployments: [CLI Overview](cli/index.md)
- Work queues & streams: [SQL Overview](sql/index.md)  
