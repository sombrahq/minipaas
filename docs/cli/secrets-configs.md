---
title: Secrets & Configs
summary: Manage Docker secrets and configs with automatic Compose patching.
---

## Overview

MiniPaaS provides high-level commands for managing Docker **secrets** and **configs**:

- Creates Swarm secrets/configs  
- Automatically updates the correct Compose file(s)  
- Patches only the services you explicitly mark with `--for`  
- Operates using multi-file service lookup based on `minipaas.yaml`

This keeps sensitive values **out of your repository**, while keeping Compose definitions **explicit and versioned**.

---

## How It Works

1. The CLI loads all Compose files listed in `minipaas.yaml`.  
2. It finds the *first* file where each `--for <service>` appears.  
3. It inserts or removes the secret/config reference into that Compose file.  
4. It creates or deletes the Docker secret/config in Swarm.  

Nothing else in the Compose file is touched.

---

## Creating Secrets

Use `secret create` to create a Swarm secret and patch Compose files:

```bash
echo postgres | minipaas secret create \
  --env dev \
  --name postgres_password \
  --for postgres \
  --for api \
  --for worker
````

This does **two things**:

1. Creates a Swarm secret named `postgres_password`.
2. Updates each service (`postgres`, `api`, `worker`) in the correct Compose file so they reference it.

Secrets are **not stored** in the Compose files — only the secret name is added.

---

## Deleting Secrets

```bash
minipaas secret delete --env dev --name postgres_password
```

This removes:

* The secret from Swarm
* All Compose references for any service that used it

---

## Pruning Unused Secrets

```bash
minipaas secret prune --env dev
```

Add `--delete` if you want the CLI to actually remove unused secrets:

```bash
minipaas secret prune --env dev --delete
```

This is safe in CI to keep the cluster clean.

---

## Creating Configs

Configs follow the same pattern but are intended for non-sensitive files.

```bash
minipaas config create \
  --env dev \
  --name app.json \
  ./configs/app.json \
  --for api \
  --for worker
```

This:

* Uploads `app.json` as a Swarm config
* Patches Compose for `api` and `worker` only
* Never modifies unrelated services

---

## Deleting Configs

```bash
minipaas config delete --env dev --name app.json
```

This removes the config and cleans references.

---

## Pruning Configs

```bash
minipaas config prune --env dev
# add --delete to remove them from Swarm
```

---

## Best Practices

* **Never commit secret values** — only commit Compose references.
* Use different secrets per environment (`dev`, `staging`, `prod`).
* Run prune operations periodically to avoid stale cluster state.
* When debugging, use Swarm introspection: `docker secret ls`, `docker config ls`.
* Re-running `create` with the same name simply replaces the value; Compose references stay intact.
