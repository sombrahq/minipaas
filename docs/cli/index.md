---
title: CLI Overview
summary: High-level overview of the MiniPaaS CLI.
---

## Overview

The MiniPaaS CLI provides a deployment and operations workflow on top of Docker Swarm while staying fully Compose-first.  
It lives in the source folder **`minipaas-cli/`** and can be used independently of the Ansible role or SQL components.

The CLI takes your existing Compose files and adds conventions for deployments, rollouts, routing, background workers, jobs, cron tasks, and secret/config management — all without introducing a control plane.

## Features

- **Compose → Swarm deployments** with reproducible builds and versioned rollouts
- **Canaries and controlled rollouts** for safe updates
- **Routes & ingress mapping** via Caddy integration
- **Workers and jobs** defined from existing Compose services
- **Cron scheduling** via swarm-cronjob
- **Secrets and configs** with automatic multi-file service lookup
- **Environment-aware execution** via `minipaas.yaml`
- **TLS helpers** for generating Docker API certificates
- **Shell helpers** for exporting DOCKER\_* context variables

## Motivation

The CLI exists to simplify deployments for teams that want:

- A Compose-first workflow that still supports multi-node clusters
- Safe, repeatable deployments without Kubernetes
- A single tool that handles secrets, configs, routing, and rollouts
- CI-friendly behavior with no external dependencies
- Minimal operational overhead

## Behavior

At a high level, the CLI:

- Reads an **environment directory** and its `minipaas.yaml` file
- Collects all referenced Compose files into a unified deployment view
- Resolves services across multiple files when modifying secrets/configs or scaffolding
- Builds and tags images consistently for deployments
- Applies rollouts and changes to Swarm services in a controlled sequence
- Integrates with Caddy, swarm-cronjob, and Docker TLS when configured

The CLI is intentionally explicit: all changes are visible in Compose files, Swarm services, or versioned artifacts.

## Start Here

- **Installation** — [CLI Installation](installation.md)
- **Configuration** — [minipaas.yaml specification](minipaas-file.md)
- **Using the CLI** — [Usage guide](usage.md)
- **Command reference** — [Commands](commands.md)

For related components:

- [Ansible Role Overview](../role/index.md)
- [SQL Overview](../sql/index.md)
