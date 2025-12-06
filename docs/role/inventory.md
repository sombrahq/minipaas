---
title: Inventory
summary: How to structure the Ansible inventory for provisioning MiniPaaS servers.
---

## Overview

The MiniPaaS Ansible role (`minipaas-role/`) uses a simple inventory layout.  
Hosts are grouped as **managers** or **workers**, and optional variables can be defined through `group_vars`, `host_vars`, or inline in the inventory.

Only four configuration variables exist:

- `docker_tls_dir`
- `minipaas_extra_ports`
- `monitoring_telegram_token`
- `monitoring_telegram_chat_id`

This keeps the provisioning layer small, predictable, and easy to maintain.

---

## Basic Inventory Structure

A minimal `inventory.ini`:

```ini
[managers]
manager1 ansible_host=1.2.3.4

[workers]
worker1 ansible_host=1.2.3.5
````

The role installs Docker, configures the system, and performs `swarm init` or `swarm join` depending on the group.

---

## Group Variables

Group variables apply to all hosts in a group.
Examples:

### `group_vars/managers.yml`

```yaml
minipaas_extra_ports:
  - 443
```

### `group_vars/workers.yml`

```yaml
minipaas_extra_ports:
  - 3000
```

---

## Host Variables

Use `host_vars/<hostname>.yml` to configure specific hosts.

Example: enabling Docker API TLS only on the manager:

```yaml
docker_tls_dir: "./certs"
```

Or enabling Telegram alerts on a specific node:

```yaml
monitoring_telegram_token: "123:ABC"
monitoring_telegram_chat_id: "999111222"
```

---

## Inline Variables in the Inventory

For small setups, variables can be set directly in `inventory.ini`:

```ini
[managers]
manager1 ansible_host=1.2.3.4 docker_tls_dir=./certs

[workers]
worker1 ansible_host=1.2.3.5 minipaas_extra_ports="[3000]"
```

Inline variables override defaults and group vars.

---

## Recommended Structure

```
inventory.ini
group_vars/
  managers.yml
  workers.yml
host_vars/
  manager1.yml
  worker1.yml
```

This approach keeps configurations organized and easy to scale.

---

## Summary

Your inventory defines:

* which hosts are managers
* which hosts are workers
* optional TLS settings
* optional firewall ports
* optional monitoring/alerting configuration

All other provisioning behavior is automatic and requires no additional settings.
