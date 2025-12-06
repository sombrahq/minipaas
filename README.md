# MiniPaaS

MiniPaaS is a lightweight toolkit for teams who want simple, repeatable infrastructure without adopting Kubernetes or hand-crafted shell scripts.  
It is built around **Docker Swarm**, **Ansible**, and **pure SQL**, with a focus on clarity, low lock-in, and practical workflows for small deployments.

MiniPaaS consists of three independent components:

- **minipaas-role/** – An Ansible role that prepares servers for Docker Swarm: Docker installation, Swarm init/join, firewall, logging, TLS, monitoring, and host-level swarm-cronjob.
- **minipaas-cli/** – A Go CLI for defining services, workers, jobs, cron tasks, building/pushing images, rolling out stacks, and managing secrets/configs.
- **minipaas-sql/** – PostgreSQL migrations implementing durable queues and streams using SQL only (dbmate-ready).

You can adopt each component separately or use them together:
- Use the **role** to prepare nodes,
- Use the **CLI** to deploy workloads,
- Use the **SQL layer** to power background workers, jobs, and event flows.

MiniPaaS is designed for:
- side projects,
- small production clusters,
- self-hosted environments,
- teams who want simplicity over control planes.

It scales from a single VM to a modest Swarm cluster with minimal cognitive load.

---

# minipaas-role

## Purpose  
A focused Ansible role that prepares plain Linux hosts to run Docker Swarm reliably and securely.

## What it sets up
- Docker CE installation
- Swarm initialization or join
- Optional Docker API TLS
- nftables firewall (allow + user-defined ports)
- syslog-ng for system and Docker logs
- fail2ban (SSH protection)
- lightweight monitoring with optional Telegram alerts
- installation of **swarm-cronjob** on the host

The role performs **host provisioning only**.  
All runtime stack deployments are handled later by the CLI.

## Minimal Inventory Example

```ini
[managers]
manager1 ansible_host=192.0.2.10 ansible_user=root

[workers]
worker1 ansible_host=192.0.2.11 ansible_user=root
````

## Role Variables

| Variable                      | Default | Purpose                                              |
| ----------------------------- | ------- | ---------------------------------------------------- |
| `docker_tls_dir`              | unset   | Directory containing TLS certificates for Docker API |
| `minipaas_extra_ports`        | `[]`    | Additional TCP ports to allow through nftables       |
| `monitoring_telegram_token`   | unset   | Token for Telegram monitoring messages               |
| `monitoring_telegram_chat_id` | unset   | Chat ID for Telegram alerts                          |

Apply the role:

```bash
ansible-playbook -i inventory.ini install.yml
```

After the role runs, the cluster is ready for deployments via the MiniPaaS CLI.

---

# minipaas-cli

## Purpose

A single binary that turns Compose definitions into Swarm rollouts and provides helpers for secrets, configs, TLS, jobs, workers, cron, and routing.

Install:

```bash
go install github.com/sombrahq/minipaas/minipaas-cli/cmd/minipaas@main
```

## Major Features

* **Service scaffolding**: exposed services, workers, jobs, cron tasks
* **Swarm automation**: build, rollout, canary releases, routing updates
* **Secrets & configs**: create, delete, prune, auto-patch Compose
* **TLS tooling**: generate server + client certificates
* **Shell environment**: load all required DOCKER_* env vars

## Command Groups

### `certs`

Generate server and client TLS certificates.

### `code`

Scaffold application behavior:

* `init`, `expose`, `worker`, `job`, `cron`

### `secret` / `config`

Create, delete, prune Docker secrets/configs with automatic Compose updates.

### `deploy`

Build images, push to registry, deploy stacks, update routing, run canaries.

### `shell`

Start a shell preloaded with TLS Docker environment variables.

---

# minipaas-sql

## Purpose

Durable background processing using PostgreSQL only.
No external brokers required.

The SQL package provides:

* **Durable queues** (enqueue, dequeue, ack)
* **Durable streams** (publish, consume with offsets)
* dbmate-ready migration files
* JSONB payloads, row locking, predictable concurrency

These SQL functions integrate naturally with MiniPaaS workers, jobs, and cron tasks.

### Applying the SQL

```bash
dbmate -m minipaas-sql/ up
```

### Example: Queue

```sql
SELECT minipaas_queue_enqueue('emails_q', '{"email":"a@b.com"}');
SELECT minipaas_queue_dequeue('emails_q', 10);
SELECT minipaas_queue_ack('emails_q', ARRAY[1]);
```

### Example: Stream

```sql
SELECT minipaas_stream_publish('activity_s', '{"action":"signup"}');
SELECT minipaas_stream_consume('activity_s', 42, 25);
```

---

# Contributing

Use conventional commits with component scopes (`feat(cli): ...`, `fix(role): ...`, `docs: ...`).
Include tests where appropriate.

---

# License

MIT License — see `LICENSE`.

