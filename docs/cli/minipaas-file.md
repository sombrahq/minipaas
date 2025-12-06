---
title: minipaas.yaml
summary: Environment configuration file for MiniPaaS, including current and planned fields.
---

## Overview

`minipaas.yaml` defines how the MiniPaaS CLI interprets an environment directory.  
It controls:

- Which Compose files belong to the environment  
- How the CLI connects to the Docker API  
- Which version tag is used for deployments  

MiniPaaS currently supports a minimal schema, but additional fields are planned.  
This page documents both:

1. **Current fields** (implemented today)  
2. **Future fields** (commented out, with notes)

---

# Current Schema (as implemented today)

```yaml
project:
  files:
    - compose.yaml
    - compose.build.yaml

api:
  host: tcp://1.2.3.4:2376   # optional
  tls: .certs                # optional (path to TLS directory)
  local: false               # if true, ignore host/tls and use local Docker

deploy:
  version: v1                # deployment version/tag
````

These are the only supported fields at the moment.

---

## Field Reference (Current)

### `project.files`

List of Compose files used by the environment.

The CLI will:

* Load them **in order**
* Resolve scaffolding modifications (`code route`, `code worker`, etc.)
* Apply secrets/configs patching to the *first* file where each service appears
* Build & deploy services defined in these files

Example:

```yaml
project:
  files:
    - compose.yaml
    - compose.build.yaml
```

---

### `api.host`

Address of the Docker API endpoint, typically TLS-secured:

```yaml
api:
  host: tcp://1.2.3.4:2376
```

If unset, the CLI uses your default Docker context.

---

### `api.tls`

Path to the directory containing:

```
ca.pem
client-cert.pem
client-key.pem
```

Example:

```yaml
api:
  tls: .certs
```

---

### `api.local`

If `true`, the CLI uses the **local Docker engine** (unix socket), ignoring `host` and `tls`.

```yaml
api:
  local: true
```

Useful for local development.

---

### `deploy.version`

The version tag to apply to built images:

```yaml
deploy:
  version: v2
```

All built images are tagged using this version.

---

# Planned / Future Fields

(These are **not supported yet**, but included here for roadmap clarity.)

```yaml
# stack: myapp
#   # NOT SUPPORTED YET
#   # Intended: allow explicit Swarm stack naming instead of deriving from env directory

# registry: registry.example.com/team
#   # NOT SUPPORTED YET
#   # Intended: automatically prefix built images with a registry path

# metadata:
#   environment: production
#   owner: backend-team
#   # NOT SUPPORTED YET
#   # Intended: CI/automation metadata container

# env:
#   APP_ENV: production
#   LOG_LEVEL: info
#   # NOT SUPPORTED YET
#   # Intended: pass env vars to build & deploy phases
```

These fields are **ignored** by the current CLI and safe to include only for documentation or future migration purposes.

---

## Example With Supported + Future Fields

```yaml
project:
  files:
    - compose.yaml
    - compose.build.yaml

api:
  host: tcp://10.0.0.5:2376
  tls: .certs
  local: false

deploy:
  version: v2025-02-01

# stack: myapp                 # future
# registry: registry.example.com/team   # future
# metadata:
#   environment: staging       # future
# env:
#   LOG_LEVEL: info            # future
```

This file is **fully valid today**, because unsupported fields remain commented.

---

## Best Practices

* Use one environment directory per deployment target (`dev/`, `staging/`, `prod/`).
* Commit `minipaas.yaml` to Git — it is part of the environment definition.
* Keep Compose files listed in `project.files` consistent across environments.
* Store TLS files in a directory like `.certs/` referenced via `api.tls`.
* Use semantic versioning, Git SHAs, or timestamp versions for `deploy.version`.
* Do not store secrets here — use `minipaas secret` / `config` instead.

