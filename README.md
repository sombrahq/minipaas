# MiniPaaS

MiniPaaS is a **middle-ground PaaS toolkit** for teams that find public-cloud platforms (AWS, Azure, GCP) too heavy and vendor-locking, yet don’t want the toil of hand-rolled servers.  
It combines three lightweight, optional components that sit on top of **Docker Swarm**:

- **minipaas-role** – An Ansible role that provisions secure Swarm nodes (TLS, firewall, monitoring, Caddy, cron, logging).
- **minipaas-cli** – A Go CLI that turns Compose files into Swarm rollouts, canaries, jobs, workers, cron tasks, and routes.
- **minipaas-sql** – PostgreSQL migrations that add durable queues and streams—no RabbitMQ/Kafka required for small apps.

**Why?**
* Deploying side projects to production quickly becomes a maze of TLS, firewalls, secrets, rollouts, and monitoring.
* Kubernetes is powerful but overkill for a cluster of “a few-hundred nodes and a few-thousand containers.”
* Vanilla Docker Compose stops at a single host; Docker Swarm gives clustering but lacks opinionated defaults.

MiniPaaS fills that gap by:

* Using **Docker Swarm** for clustering (simpler than Kubernetes, bigger than Compose).
* Shipping drop-in **Ansible** and **CLI** automation so you don’t have to write shell scripts.
* Relying on **open-source defaults** (Caddy, nftables, syslog-ng, telegram-alerts) to avoid SaaS lock-in.
* Letting you adopt each part independently—provision nodes with the role, or just use the CLI, or only the SQL framework.

Start on a single VM, grow to a modest Swarm, and postpone a Kubernetes migration until it’s truly justified—all while keeping infra cost and cognitive load in check.

---

## minipaas-role

### Purpose
An Ansible role that turns a vanilla Linux box into a secure Docker Swarm node (manager or worker) with extra runtime services.

### Features

* Installs Docker CE and initialises or joins Swarm (TLS enabled).
* Inventory-driven manager / worker roles.
* Installs and configures:
    * **swarm-cronjob** – cron in Swarm.
    * **telegram-bash-system-monitoring** – lightweight node metrics + Telegram alerts.
    * **fail2ban** – basic SSH/Docker API brute-force protection.
    * **syslog-ng** – error logs forwarded to Telegram.
    * **nftables** – simple host firewall.
    * Optional **Caddy** reverse proxy for HTTPS + automatic certs.

---

### Environment variables (`.env` file)


| Variable                | Purpose / Example value                                   |
|-------------------------|-----------------------------------------------------------|
| `TELEGRAM_BOT_TOKEN`    | Bot token for Telegram alerts (`123456:ABC-DEF…`)         |
| `TELEGRAM_CHAT_ID`      | Chat or channel ID that will receive notifications        |

These env vars are loaded by `ansible.builtin.env` look-ups inside the role.

---

### Role variables

| Variable                      | Default value         | Purpose / Notes                                                                                                                         |
|-------------------------------|-----------------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| `docker_tls_dir`              | `.certs`              | Directory (on the **control-host**) that contains `ca.pem`, `server-cert.pem`, and `server-key.pem`. When present, the role enables TLS-authenticated access to the Docker API. |
| `swarm_cronjob_version`       | `1.14.0`              | Version of **swarm-cronjob** to download and install on the node designated as the cron host.                                           |
| `telegram_bot_token`          | `$TELEGRAM_BOT_TOKEN` | Telegram Bot token used by syslog-ng and the monitoring script to send alerts. Usually injected via `.env`.                             |
| `telegram_chat_id`            | `$TELEGRAM_CHAT_ID`                   | Chat or channel ID that receives Telegram alerts.                                                                                       |
| `minipaas_runtime_enabled`    | `true`                | When `true`, the role deploys the optional **MiniPaaS runtime stack** (Caddy, registry, etc.) on the Swarm.                             |
| `minipaas_caddy_enabled`      | `true`                | Enables or disables deployment of the Caddy reverse-proxy service.                                                                      |
| `caddy_image_tag`             | `2.9.1-alpine`        | Docker image tag used for the Caddy service.                                                                                            |
| `minipaas_registry_enabled`   | `true`                | Deploys a private Docker registry inside the Swarm when set to `true`.                                                                  |
| `minipaas_registry_image_tag` | `2.9.1-alpine`        | Docker image tag used for the private registry service.                                                                                 |


Override any of these in **`group_vars/all.yml`** (or host-specific var files).

---

### Inventory example (`inventory.ini`)

```ini
[managers]
manager1 ansible_host=192.0.2.10 ansible_user=root

[workers]
worker1  ansible_host=192.0.2.11 ansible_user=root
worker2  ansible_host=192.0.2.12 ansible_user=root

```

---

### Minimal playbook (`install.yml`)

```yaml
---
- hosts: all
  become: yes

  roles:
    - minipaas-role            # <-- relative or Galaxy install name
```

---

### Running the role

```bash
# 1️⃣ make sure your .env and inventory.ini files are ready
# 2️⃣ execute the playbook
ansible-playbook -i inventory.ini install.yml
```

**What happens next?**

* Docker is installed and initialised (`docker swarm init` on the first manager).
* Certificates are copied to `/etc/docker` and `daemon.json` is patched for TLS.
* Swarm-Cronjob, syslog-ng + Telegram, fail2ban, nftables, and Caddy are installed + enabled.
* Managers get the “manager” label and workers join automatically with the stored token.
* Your fleet is now ready for **MiniPaaS CLI**, CapRover, Portainer, or any Swarm workload.

> **Tip:** use `--tags docker` or `--tags monitoring` to run only a subset of the role when iterating.

### Alternatives

* **CapRover** – includes its own provisioning, but Swarm nodes can be prepped with minipaas-role.
* **Portainer** – UI for Swarm; does not provision the host.
* **Swarmpit** – metrics dashboard/UI atop Swarm but lacks automation provisioning.

---

## minipaas-cli

### Purpose
A single binary that drives Docker Swarm operations over TLS and augments Compose files with Swarm best practices.

### Features
* Secrets & Configs: create, delete, prune, auto-inject into Compose.
* Code generators: exposed services, workers, cron jobs, one-off jobs.
* Deployments: `deploy build`, `deploy rollout`, `deploy canary`, `deploy routing`.
* TLS helper commands (`certs server`, `certs client`).

### Required commands

* ``docker``
* ``openssl``
* ``cp``


### certs — TLS Certificate Management
| Command | What it does |
|---------|--------------|
| `minipaas certs server`  | Generates **server-side** TLS certificates (CA, key, cert) for the Docker API on Swarm nodes. |
| `minipaas certs client`  | Generates **client-side** certificates so the CLI (or CI) can authenticate against the TLS-secured Docker API. |

### code — Environment & Service Scaffolding
| Command | What it does |
|---------|--------------|
| `minipaas code init`    | Creates a new *MiniPaaS environment* directory with `minipaas.yaml`, compose overrides, and Caddy config. |
| `minipaas code expose`  | Adds Caddy routing and Swarm labels to expose a service at a public URL. |
| `minipaas code job`     | Marks one or more services as **one-off jobs** (run once, exit). |
| `minipaas code worker`  | Converts services into **background workers** (long-running, no external port). |
| `minipaas code cron`    | Schedules services to run on a **cron** expression using `swarm-cronjob`. |

### secret — Docker Secret Management
| Command | What it does |
|---------|--------------|
| `minipaas secret create` | Creates a Docker secret from a file or STDIN, hashes the content for versioning, and patches compose files. |
| `minipaas secret delete` | Removes a secret from the Swarm and updates compose files to drop references. |
| `minipaas secret prune`  | Lists (or deletes with `--delete`) secrets that are no longer attached to any running service. |

### config — Docker Config Management
| Command | What it does |
|---------|--------------|
| `minipaas config create` | Same workflow as `secret create`, but for *Docker configs* (non-sensitive files). |
| `minipaas config delete` | Deletes a config and cleans references from compose/service definitions. |
| `minipaas config prune`  | Lists or deletes unused configs across the Swarm. |

### deploy — Build & Release Automation
| Command | What it does |
|---------|--------------|
| `minipaas deploy build`    | Runs `docker compose build --push` with proper tags for the private registry. |
| `minipaas deploy rollout`  | Performs a zero-downtime **rolling update** of the Swarm stack. |
| `minipaas deploy canary`   | Adds *N* extra replicas (default 1) for a **canary release**, then scales back when validated. |
| `minipaas deploy routing`  | Pushes the latest Caddy HTTP server configuration into the runtime (via Caddy Admin API). |

### shell — Convenience Helper
| Command | What it does |
|---------|--------------|
| `minipaas shell` | Spawns an interactive shell with all required `DOCKER_*` environment variables exported. |

---

> **Tip**   
> Run `minipaas --help` or `minipaas <subcommand> --help` to see full flag reference, examples, and advanced options.


### Alternatives
* **Kamal** – SSH + plain Docker (no Swarm).
* **CapRover CLI** – tightly coupled with CapRover dashboards.
* **Docker Compose CLI** – good for single host; no Swarm rollouts/canaries.

---

## minipaas-sql

### Purpose
`minipaas-sql` brings **durable queues** and **append-only streams** to PostgreSQL by means of PL/pgSQL functions.  
It is ideal for teams following **Hexagonal / Clean Architecture** because the message-passing boundary is abstracted behind SQL calls—if you later swap Postgres for RabbitMQ, Kafka, or anything else, only the adapter layer changes, not your core domain code.

### Features
| Category | Functions | Highlights |
|----------|-----------|------------|
| **Queues** | `minipaas_queue_create` · `minipaas_queue_enqueue` · `minipaas_queue_dequeue` · `minipaas_queue_ack` | Dynamic table creation, batch dequeue with row locking, `pg_notify` on enqueue. |
| **Streams** | `minipaas_stream_create` · `minipaas_stream_publish` · `minipaas_stream_consume` | Append-only event tables, consumer checkpoints, `pg_notify` fan-out. |
| **Safety** | - | Table-name regex validation and batch-size checks prevent accidental SQL injection. |
| **DevOps** | - | Ships as plain SQL so you can apply with any migration tool (`dbmate`, Flyway, Sqitch, Liquibase…). |

### Usage

#### Apply migrations

```bash
dbmate -u "postgres://user:pass@host/db" -m minipaas-sql/ up
```

#### Working with a queue

```sql
-- 1. Create the queue table
SELECT minipaas_queue_create('emails_q');

-- 2. Enqueue a message (triggers LISTEN/NOTIFY on channel "emails_q")
SELECT * FROM minipaas_queue_enqueue('emails_q', '{"email":"a@b.com"}'::jsonb);

-- 3. Dequeue 10 jobs for processing
SELECT * FROM minipaas_queue_dequeue('emails_q', 10);

-- 4. Acknowledge completion
SELECT minipaas_queue_ack('emails_q', ARRAY[1,2,3]);
```

#### Working with a stream

```sql
-- 1. Create the stream
SELECT minipaas_stream_create('activity_s');

-- 2. Publish an event
SELECT * FROM minipaas_stream_publish('activity_s', '{"action":"signup"}'::jsonb);

-- 3. Consume events after id 42 (batch of 25)
SELECT * FROM minipaas_stream_consume('activity_s', 42, 25);
```

#### Listening for new items

```sql
-- Receive NOTIFY payloads whenever a new job/event arrives
LISTEN "emails_q";
LISTEN "activity_s";
```

### Alternatives
* **RabbitMQ (AMQP)** – Advanced routing, per-message TTL, pluggable extensions.
* **Apache Kafka** – High-throughput commit log, partitions, rich stream-processing ecosystem.
* **Redis Streams** – In-memory speed and simplicity; persistence is optional.
* **pgmq** – Another Postgres queue library with visibility timeouts and DLQs.
* **AWS SQS / Google PubSub** – Fully managed queues when you are already invested in those clouds.


`minipaas-sql` lets you start small—staying inside Postgres—and graduate to any specialized message broker later without rewriting the heart of your application.

---

## Contributing

Pull requests and issues are welcome! Please follow conventional commit messages and include tests where reasonable.

---

## License

MiniPaaS is released under the MIT License. See the `LICENSE` file for details.
