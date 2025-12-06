---
title: Commands
summary: Reference for all MiniPaaS CLI commands.
---

# Overview

MiniPaaS CLI commands are organized into functional groups:

- **code** → modify compose files (routing, workers, jobs, cron, init)  
- **secret / config** → manage Swarm secrets and configs with automatic compose updates  
- **deploy** → build images, roll out updates, and apply routing  
- **certs** → generate TLS certificates for Docker API access  
- **shell** → open a Docker-ready shell for local or remote contexts  

All commands operate relative to an environment directory (`--env <dir>`) that contains `minipaas.yaml`.

---

# Global Flags

- `--env <dir>` — environment directory containing `minipaas.yaml`  
- `--verbose` — detailed output  
- `--files` (for `code init`) — list of compose files to wire into the environment  
- `--for <service>` — service(s) to attach secrets/configs to  
- `--name <name>` — name of secret/config  

---

# `minipaas code` — Compose-aware scaffolding

These commands **modify your compose files in place**.

### Initialize an environment

```bash
minipaas code init --env dev \
  --files compose.yaml compose.build.yaml
````

Creates/updates:

* `dev/minipaas.yaml`
* compose-file mapping inside the environment

---

### Route a service via Caddy

```bash
minipaas code route --env dev \
  http://localhost:8000 \
  api:8080
```

Arguments:

1. Public URL
2. `<service:port>` pair inside the Swarm stack

Adds routing labels to the appropriate compose file.

---

### Mark background worker services

```bash
minipaas code worker --env dev worker emails_worker
```

Updates each service definition to behave as a long-running worker.

---

### Mark a job service (run-once)

```bash
minipaas code job --env dev migrate
```

Job services run once and exit. Execution is triggered via deploy workflows or CI.

---

### Mark a cron-managed service

```bash
minipaas code cron --env dev cleanup
```

Makes the service runnable through `swarm-cronjob`. Schedule is configured in compose/env configuration.

---

# `minipaas secret` — Swarm secrets (multi-file compose patching)

### Create a secret

```bash
echo postgres | minipaas secret create \
  --env dev \
  --name postgres_password \
  --for postgres --for api --for migrate
```

Creates a Swarm secret and patches the compose file where each service is defined.

---

### Delete a secret

```bash
minipaas secret delete --env dev --name postgres_password
```

Removes the secret and cleans references from affected compose files.

---

### Prune unused secrets

```bash
minipaas secret prune --env dev
# add --delete to remove unused secrets
```

---

# `minipaas config` — Swarm configs (same model as secrets)

### Create a config

```bash
minipaas config create \
  --env dev \
  --name app.json \
  ./configs/app.json \
  --for api --for worker
```

### Delete a config

```bash
minipaas config delete --env dev --name app.json
```

### Prune unused configs

```bash
minipaas config prune --env dev
# add --delete to remove unused configs
```

---

# `minipaas deploy` — Build, rollout, routing

These commands operate on the environment’s compose files + stack config.

### Build images

```bash
minipaas deploy build --env dev
minipaas deploy build --verbose --env dev
```

Builds and tags all images referenced in `minipaas.yaml`.

---

### Rollout an update

```bash
minipaas deploy rollout --env dev
minipaas deploy rollout --verbose --env dev
```

Applies controlled updates to services in the Swarm stack.

---

### Apply routing configuration (Caddy)

```bash
minipaas deploy routing --env dev
```

Updates routing files/services after `code route` or Caddy config changes.

---

### Canary (optional depending on version)

```bash
minipaas deploy canary --env dev --replicas 1 api worker
```

Temporarily runs additional replicas for testing.

---

# `minipaas certs` — TLS for Docker API

Generate certificates for Swarm manager (server) or CLI/CI (client).

### Generate all certificates

```bash
minipaas certs all --out .certs
```

### Generate server-only certificates

```bash
minipaas certs server --out .certs
```

### Generate client-only certificates

```bash
minipaas certs client --out .certs
```

---

# `minipaas shell` — Docker-ready shell

Open a shell with `DOCKER_HOST`,`DOCKER_TLS_VERIFY`, `DOCKER_CERT_PATH`, etc. properly set:

```bash
minipaas shell --env dev
```

Useful for debugging or CI environments.

---

# Notes

* All compose modifications are **idempotent**.
* Secrets/configs patch only the first compose file where each service appears.
* Deploy commands expect the remote Docker API (TLS optional) to already be reachable.
* The CLI never overwrites unrelated sections inside compose files.
