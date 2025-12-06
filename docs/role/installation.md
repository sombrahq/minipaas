---
title: Installation
summary: Install and run the MiniPaaS Ansible role to prepare servers for a minimal Docker Swarm cluster.
---

## Overview

The MiniPaaS role provisions servers so they are ready to participate in a **clean, empty Docker Swarm cluster**.  
It does not deploy any runtime services. The MiniPaaS CLI handles all Swarm workloads later.

The role installs:

- Docker
- Swarm (init + join)
- nftables firewall
- optional Docker API TLS
- syslog-ng
- optional Fail2Ban
- optional lightweight monitoring & Telegram alerts
- `swarm-cronjob` installed

---

## Requirements

You need:

- Linux hosts with SSH access
- Python 3 and Ansible installed locally
- An inventory file defining `managers` and `workers`
- (Optional) TLS certificates for Docker API access
- (Optional) tokens/IDs for monitoring and alerting

Supported operating systems:

- Debian / Ubuntu family
- Any modern Linux using systemd and nftables

---

## Install the Role Locally

Clone the repository:

```bash
git clone https://github.com/sombrahq/minipaas.git
cd minipaas
````

Install Ansible Galaxy requirements (only if present):

```bash
ansible-galaxy install -r requirements.yml
```

---

## Prepare an Inventory

A minimal structure:

```ini
[managers]
manager1 ansible_host=1.2.3.4

[workers]
worker1 ansible_host=1.2.3.5
```

The role automatically:

* initializes Swarm on the first manager
* uses join-tokens to attach additional managers/workers

More details: **[Inventory](inventory.md)**

---

## Optional: Configure TLS for Remote Docker API

To enable secure remote access for CI or the CLI:

1. Generate certificates using the MiniPaaS CLI or your own PKI
2. Place them in a directory (preferably stored via Ansible Vault)
3. Set:

```yaml
docker_tls_dir: "./certs"
```

Providing these files enables the Docker API at:

```
tcp://0.0.0.0:2376 (TLS required)
```

If omitted → Docker is Unix-socket only.

---

## Run the Role

Provision the entire cluster:

```bash
ansible-playbook -i inventory.ini main.yml
```

The role performs:

* Docker installation
* Swarm init on the manager
* Swarm join on workers
* nftables firewall configuration
* syslog-ng configuration
* optional Fail2Ban
* optional monitoring script + Telegram alerts
* installation of `swarm-cronjob` on the host

The playbook is **idempotent** — you can apply it repeatedly to converge state safely.

---

## Verifying Installation

### Check Docker

```bash
docker info
```

### Check that Swarm is active

```bash
docker node ls
```

### Check that `swarm-cronjob` is installed on the host

```bash
systemctl status swarm-cronjob
```

(or depending on your setup, a cron entry may be present)

### Check Docker API TLS (if enabled)

```bash
docker --tlsverify \
  --tlscacert=ca.pem \
  --tlscert=client-cert.pem \
  --tlskey=client-key.pem \
  -H tcp://<manager-ip>:2376 info
```

---

## What Happens Next?

Once the role finishes:

1. Machines are hardened and configured
2. The infrastructure is ready for application/runtime deployment
3. The MiniPaaS CLI will create networks, routing, services, secrets, configs, workers, jobs, cron services, etc.

To proceed:

* Move to the CLI → **[MiniPaaS CLI Overview](../cli/index.md)**

---

## Summary

Use the MiniPaaS Ansible role to:

* prepare machines
* bootstrap Swarm
* secure the environment
* enable optional remote API access
* optionally install host-level cron orchestration
* standardize logs and firewall
