---
title: Swarm
summary: How the MiniPaaS role initializes and joins a Docker Swarm cluster.
---

## Overview

The MiniPaaS role provisions **Docker Swarm**, but does so in a minimal, infrastructure-only way:

- It installs Docker on each host.  
- It initializes a Swarm manager on the primary node.  
- It joins all worker nodes to the cluster.  
- It does **not** deploy any services, overlay networks, or runtime stack components.  
- It installs `swarm-cronjob` **on the host**, not inside the Swarm cluster.

The result is a **clean, empty Swarm cluster** ready for the MiniPaaS CLI (or your own tooling) to deploy applications and runtime components.

This separation keeps the role focused on repeatable host provisioning and leaves orchestration behavior to the CLI.

---

## What the Role Does

### 1. Installs Docker
All required packages are installed so the host can run containers and participate in a Swarm cluster.

### 2. Initializes or Joins Swarm
- The first manager performs `docker swarm init`.
- Additional managers and workers join using the appropriate tokens.
- The cluster is left in a stable, but **empty**, state.

No services, stacks, or networks are created at this stage.

### 3. Configures Swarm Ports
The role ensures nftables allows required Swarm traffic:

- TCP 2377 (Swarm management)
- TCP/UDP 7946 (gossip)
- UDP 4789 (overlay networking)

### 4. Installs swarm-cronjob on the Host
If enabled:

- `swarm-cronjob` is installed **on the manager node**, not inside Swarm.
- It runs as a system process (e.g. via systemd or cron, depending on your setup).
- It communicates with the Docker Engine via the API (Unix socket or TLS).
- It triggers Swarm tasks according to cron expressions found in service labels.

This enables cron-driven Swarm tasks without requiring a running daemon inside the cluster.

---

## Cluster State After Provisioning

After the role finishes, your infrastructure looks like this:

### Managers
- Docker installed  
- Swarm initialized  
- Docker API optionally secured with TLS  
- nftables firewall active  
- syslog-ng, monitoring, and fail2ban (if enabled)  
- `swarm-cronjob` installed on the host (optional)  

### Workers
- Docker installed  
- Joined to the Swarm  
- nftables rules applied  
- syslog-ng (if enabled)  

The cluster is operational but intentionally empty.

---

## Deploying Into the Cluster

Once the Swarm is created by the role, the MiniPaaS CLI handles:

- creating overlay networks  
- deploying Caddy ingress  
- deploying Postgres (if needed)  
- deploying application stacks  
- routing updates  
- worker/job/cron scaffolding  
- secrets and configs  
- rolling updates

This ensures a clean split between:

- **Infrastructure provisioning** → role  
- **Application + runtime orchestration** → CLI

---

## Best Practices

- Use a single manager unless you need HA; small clusters are simpler.  
- Keep the Swarm empty until the CLI deploys runtime components.  
- Do not manually create networks or services unless you know what you're doing — let the CLI maintain consistency.  
- If using TLS, allow only secure Docker API access (restrict port `2376` externally).  
- Treat `swarm-cronjob` on the host as part of the infrastructure layer — it will run cron-triggered tasks for any service deployed by the CLI.

---

## Summary

The MiniPaaS role prepares servers for MiniPaaS by:

- installing Docker  
- creating a clean Swarm cluster  
- configuring firewall, logging, and TLS  
- optionally installing a host-level cron scheduler for Swarm  

Everything else — routing, databases, workers, jobs, cron services, networks, stacks — is handled by the **MiniPaaS CLI**, not the role.

