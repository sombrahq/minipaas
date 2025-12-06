---
title: Usage
summary: How to use the MiniPaaS role to provision and manage Swarm nodes.
---

## Overview

This guide shows how to **run**, **re-run**, and **evolve** a MiniPaaS-powered Swarm cluster after installation.

The MiniPaaS role is designed to be:

- fully reproducible
- idempotent (safe to run repeatedly)
- easy to extend
- easy to undo
- suitable for small Swarm clusters and side projects

You describe hosts in your `inventory.ini`, adjust variables in `group_vars` / `host_vars`, and apply the role using
Ansible.

---

## Running the Role

After preparing your inventory and configuration:

```bash
ansible-playbook -i inventory.ini main.yml
````

This performs:

* Docker installation
* Swarm initialization or joining
* Firewall setup (nftables)
* syslog-ng installation
* Fail2Ban (optional)
* swarm-cronjob installation
* Caddy + Postgres runtime stack deployment (optional)
* Docker API TLS (optional)

All steps are **idempotent**: re-running the playbook will converge the state without breaking existing deployments.

---

## Re-running the Role

You should re-run the playbook when you:

* add new nodes
* rotate TLS certificates
* modify firewall rules
* enable/disable components (Fail2Ban, syslog, runtime stack, cronjob)
* update monitoring or alerting settings
* change the overlay network configuration
* want to verify cluster health after changes

Example:

```bash
ansible-playbook -i inventory.ini main.yml
```

Re-running the role is safe even in production-like environments, as long as you understand which components may
restart (Docker, Caddy, etc.).

---

## Adding a New Manager

1. Add the new host to the `managers` group
2. Set its variables in `host_vars/<name>.yml`
3. Re-run the role:

```bash
ansible-playbook -i inventory.ini main.yml
```

The role will:

* install Docker
* join the Swarm as a manager
* apply firewall rules
* configure monitoring/logging

Note: multiple managers increase operational overhead; MiniPaaS aims for **small, simple clusters**.

---

## Adding a Worker Node

1. Add the host to `[workers]`
2. (Optional) add `enable_firewall: true` or other overrides
3. Apply the role:

```bash
ansible-playbook -i inventory.ini main.yml
```

The worker will join the cluster automatically using the stored join token.

---

## Enabling TLS for Remote Docker Access

If you want CI/CD or the MiniPaaS CLI to connect remotely:

1. Generate certificates (e.g. via the CLI)
2. Place them on the manager (preferably via Ansible Vault)
3. Configure:

```yaml
docker_api_tls_dir: "./certs"
```

4. Run:

```bash
ansible-playbook -i inventory.ini main.yml
```

The Docker API will be exposed on:

```
tcp://0.0.0.0:2376 (TLS required)
```

---

## Managing the Runtime Stack

The built-in runtime stack includes:

* Caddy (ingress)
* Postgres (optional)
* Overlay network creation

To disable deployment of these components:

```yaml
deploy_runtime_stack: false
```

To disable specific components:

```yaml
deploy_caddy: false
deploy_postgres: false
```

Re-run the role to apply changes.

---

## Updating Firewall Rules

Modify allowed ports:

```yaml
firewall_allowed_tcp_ports:
  - 22
  - 80
  - 443
  - 3000
```

Then apply:

```bash
ansible-playbook -i inventory.ini main.yml
```

The nftables rules will be regenerated and applied safely.

---

## Updating Monitoring or Alerts

Enable system monitoring:

```yaml
enable_monitoring_script: true
```

Add Telegram alerts:

```yaml
monitoring_telegram_token: "123:ABC"
monitoring_telegram_chat_id: "999111222"
```

Apply changes with:

```bash
ansible-playbook -i inventory.ini main.yml
```

---

## Removing Components

Everything installed by the role is:

* standard Debian/Ubuntu packages
* standard Docker stack files
* systemd units and configuration files

To remove components:

* disable variables in `group_vars` / `host_vars`
* re-run the role
* or manually uninstall the packages / stacks

Because MiniPaaS is **low lock-in**, removal is safe and predictable.

---

## Cluster Lifecycle Examples

### Redeploy TLS

```bash
# Replace certs
ansible-playbook -i inventory.ini main.yml
```

### Rebuild firewall after opening a new port

```yaml
firewall_allowed_tcp_ports:
  - 22
  - 443
  - 8080
```

```bash
ansible-playbook -i inventory.ini main.yml
```

### Reinstall monitoring script

```yaml
enable_monitoring_script: true
```

```bash
ansible-playbook -i inventory.ini main.yml
```

### Reset a worker node

1. Reinstall OS or clean Docker
2. Add node back to `[workers]`
3. Run the playbook
4. Worker rejoins automatically

---

## Summary

Using the MiniPaaS role consists of:

* Writing a clear inventory
* Setting per-group or per-host variables
* Running Ansible to converge the cluster
* Re-running to apply changes safely

The role provides a **minimal, transparent, reproducible Swarm environment**, ideal for small deployments and teams that
value low operational overhead.
