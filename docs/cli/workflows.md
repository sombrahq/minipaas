---
title: Workflows
summary: Build, rollout, canary, and routing operations.
---

MiniPaaS provides a few high‑level workflows to operate your stack on Swarm.

## Build images

Build tagged images for your services using the compose build context.

```bash
minipaas deploy build --env dev
```

## Rolling update

Deploy a zero‑downtime rollout of the stack.

```bash
minipaas deploy rollout --env dev
```

## Canary

Add temporary replicas using the currently deployed image for a service.

```bash
minipaas deploy canary --env dev --replicas 1 api worker
```

## Routing

Push updated Caddy config to the runtime.

```bash
minipaas deploy routing --env dev
```
