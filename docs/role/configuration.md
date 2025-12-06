---
title: Configuration
summary: Variables used to configure the MiniPaaS Ansible role.
---

## Overview

The `minipaas-role/` exposes a small set of configuration variables to control:

- Docker installation  
- Swarm initialization and joining  
- TLS for the Docker API  
- Caddy + Postgres base stack deployment  
- Firewall behavior  
- Monitoring and logging  
- Cron execution via swarm-cronjob  

This page documents the key variables, where they are defined, and how to override them.

You can set these variables in:

- `group_vars/managers.yml`  
- `group_vars/workers.yml`  
- `host_vars/<hostname>.yml`  
- or directly in `inventory.ini`  

---

# Core Variables

### `docker_api_tls_dir`  
Directory containing:

```

ca.pem
server-cert.pem
server-key.pem

````

If set, the role configures Docker to expose a **TLS-secured remote API** (`tcp://0.0.0.0:2376`).

Example:

```yaml
docker_api_tls_dir: "./certs"
````

If omitted → Docker API is **not** exposed.

---

### `minipaas_overlay_network`

Name of the overlay network created for application deployments.

```yaml
minipaas_overlay_network: "minipaas_net"
```

Services deployed via the MiniPaaS CLI are typically attached to this network.

---

### `deploy_runtime_stack`

Whether to deploy the built-in Caddy + Postgres runtime stack.

```yaml
deploy_runtime_stack: true
```

If false → the role provisions the nodes but does **not** deploy any base services.

---

# Swarm Setup

### `swarm_manager`

Boolean indicating whether a node should initialize the Swarm.

Managers should set:

```yaml
swarm_manager: true
```

Workers should set:

```yaml
swarm_manager: false
```

### `swarm_manager_ip`

Required on managers for cluster initialization.

```yaml
swarm_manager_ip: "1.2.3.4"
```

Workers use the manager’s advertised address to join.

---

### `swarm_join_token`

Automatically generated on the first manager; workers use it to join.

Most users do **not** set this manually — the role handles it.

---

# Firewall Configuration

### `enable_firewall`

Enable nftables with a default-deny policy.

```yaml
enable_firewall: true
```

### `firewall_allowed_tcp_ports`

List of extra ports to allow.

By default, the role opens required Swarm ports automatically.

Example:

```yaml
firewall_allowed_tcp_ports:
  - 22     # SSH
  - 80     # HTTP
  - 443    # HTTPS
```

If you run extra services outside Caddy, add them here.

---

# Logging & Monitoring

### `enable_syslog`

Enable syslog-ng installation and configuration.

```yaml
enable_syslog: true
```

### `enable_fail2ban`

Install basic Fail2Ban protection.

```yaml
enable_fail2ban: true
```

### `enable_monitoring_script`

Install the lightweight system monitoring script that periodically summarizes CPU, memory, disk, and Swarm status.

```yaml
enable_monitoring_script: true
```

### `monitoring_telegram_token` / `monitoring_telegram_chat_id`

Optional. If set, monitoring alerts flow to a Telegram bot.

```yaml
monitoring_telegram_token: "123:ABC"
monitoring_telegram_chat_id: "987654321"
```

If not set → Telegram alerts are disabled.

---

# Cron / Scheduled Jobs

### `enable_swarm_cronjob`

Deploy `swarm-cronjob` on the main manager.

```yaml
enable_swarm_cronjob: true
```

MiniPaaS CLI's `code cron` depends on this component.

---

# Caddy Ingress Runtime

### `caddy_email`

Email used for ACME/Let’s Encrypt when TLS is enabled in Caddy.

```yaml
caddy_email: "admin@example.com"
```

If left blank, Let’s Encrypt operations may be disabled or skipped depending on your Caddy configuration.

### `deploy_caddy`

Whether to include Caddy in the runtime stack.

```yaml
deploy_caddy: true
```

---

# Postgres Runtime

### `deploy_postgres`

Default Postgres deployment behavior.

```yaml
deploy_postgres: true
```

### `postgres_password`

Password for the default Postgres instance deployed by the role.

**Important:**
This is *not* meant for long-term sensitive environments.
For production setups, consider managing credentials outside the role.

---

# Paths & Templates

### `runtime_stack_template`

Path to the Jinja template for the base Swarm stack file.

```yaml
runtime_stack_template: "templates/swarm-stack.yml.j2"
```

Users normally do not change this unless customizing the runtime.

---

# Example Configuration

```yaml
swarm_manager: true
swarm_manager_ip: "1.2.3.4"

docker_api_tls_dir: "./certs"

enable_firewall: true
firewall_allowed_tcp_ports:
  - 22
  - 80
  - 443

enable_syslog: true
enable_fail2ban: true
enable_monitoring_script: true

enable_swarm_cronjob: true

deploy_runtime_stack: true
deploy_caddy: true
deploy_postgres: true

caddy_email: "admin@example.com"
postgres_password: "example123"
```

---

# Best Practices

* Keep your inventory simple: one manager group, one worker group.
* Don’t expose Docker API TLS externally unless required.
* Store TLS certs and sensitive variables in **Ansible Vault**.
* For production Postgres setups, deploy your own cluster instead of relying on the default runtime stack.
* Re-run the role safely — it is idempotent.
* Keep configs in Git for full reproducibility.

