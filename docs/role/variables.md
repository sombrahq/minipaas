---
title: Variables
summary: Configuration variables supported by the MiniPaaS Ansible role.
---

## Overview

The MiniPaaS Ansible role (`minipaas-role/`) exposes a small set of variables that control:

- Docker API TLS
- Additional allowed TCP ports
- Telegram-based monitoring alerts

These variables can be defined in:

- `inventory.ini`
- `group_vars/`
- `host_vars/`

---

# Variables

| Variable                      | Purpose                                            | Default          |
|-------------------------------|----------------------------------------------------|------------------|
| `docker_tls_dir`              | Directory containing TLS certificates              | unset (disabled) |
| `minipaas_extra_ports`        | Additional TCP ports to allow through the firewall | `[]`             |
| `monitoring_telegram_token`   | Token for sending monitoring messages to Telegram  | unset            |
| `monitoring_telegram_chat_id` | Telegram chat ID for monitoring output             | unset            |

---

## Variable Details

### API TLS

Specifies a directory with the certificates needed to enable Docker’s remote TLS API.

Expected files inside this directory:

- `ca.pem`
- `server-cert.pem`
- `server-key.pem`

If this variable is not set, the Docker API remains accessible only via the Unix socket.

---

### Firewall

A list of extra TCP ports to allow through the firewall.

Example:

```yaml
minipaas_extra_ports:
  - 443
  - 3000
````

---

### Telegram Monitoring

Enable Telegram alerts from the host’s monitoring script.

Example:

```yaml
monitoring_telegram_token: "123:ABC"
monitoring_telegram_chat_id: "999111222"
```

If either is unset, Telegram notifications remain disabled.

---

## Summary

These four variables provide optional configurations for TLS, firewall customization, and monitoring.
All other server provisioning tasks use built-in defaults and require no additional settings.
