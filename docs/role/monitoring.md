---
title: Monitoring
summary: Lightweight monitoring and alerting features provided by the MiniPaaS role.
---

## Overview

The MiniPaaS role includes **lightweight monitoring**, designed for small clusters, homelabs, and side projects.  
It does **not** install Prometheus, Grafana, Loki, ELK, or any heavy monitoring suite.

Instead, it provides:

- a minimal **system monitoring script** (periodic report)  
- **syslog-ng** for log collection  
- **optional Telegram alerts** for basic notifications  
- simple Swarm health summaries  

These tools give enough visibility for small deployments without committing to a full observability stack.

---

## Syslog-ng

If enabled:

```yaml
enable_syslog: true
````

then the role installs **syslog-ng** and configures:

* system logs
* Docker service logs
* Swarm service logs (via journald)

Syslog-ng acts as a unifying log sink on each node.
You can use standard Linux tooling (e.g., `journalctl`, `grep`, `awk`) or forward logs elsewhere if needed.

No external log storage (Loki, Elasticsearch, etc.) is installed.

---

## Monitoring Script

If enabled:

```yaml
enable_monitoring_script: true
```

the role installs a small monitoring script that periodically reports:

* CPU usage
* Memory usage
* Disk space
* Docker status
* Swarm node state
* Number of running services
* Basic system health indicators

The script generates a compact summary useful for:

* homelabs
* small VMs
* low-cost cloud instances
* environments where full observability stacks are overkill

The script does **not** keep historical metrics; it provides snapshots only.

---

## Telegram Alerts (Optional)

If you configure:

```yaml
monitoring_telegram_token: "123:ABC"
monitoring_telegram_chat_id: "999111222"
```

the monitoring script will send its periodic report to Telegram.

This provides a simple “heartbeat-style” visibility:

* machine alive/dead
* Swarm healthy/unhealthy
* disk nearing capacity
* service count changes
* CPU/memory spikes

If these variables are not set → Telegram notifications are disabled.

---

## Cadence

The monitoring script runs via cron (or systemd timer on some systems), depending on your configuration.
Default cadence: **every few minutes** (implementation-specific).

You can modify the schedule by editing the cron file deployed to the server.

---

## Fail2Ban Alerts (Optional)

If `enable_fail2ban: true`:

* Fail2Ban monitors SSH for brute-force attempts
* Alerts can be forwarded to Telegram if monitoring variables are set
* Logs flow into syslog-ng for visibility

This setup is intentionally minimal — enough for a low-risk server, not a replacement for full intrusion detection.

---

## Removing Monitoring

All monitoring components are:

* standard packages (syslog-ng, fail2ban)
* standard cron entries or systemd timers
* simple scripts installed by the role

Removing them is as easy as:

* uninstalling packages
* deleting the monitoring script
* disabling cron entries
* re-running the role with monitoring disabled

There is **no lock-in**.

---

## When to Add a Real Monitoring Stack

MiniPaaS monitoring is intentionally minimal.
Consider a full stack if you need:

* historical metrics
* dashboards
* service-level alerting
* application performance visibility
* distributed tracing
* structured log pipelines
* multi-node log aggregation

Suggested add-ons (not provided by MiniPaaS):

* Prometheus + Grafana
* Loki + Promtail
* ELK or OpenSearch stack
* Netdata
* Datadog / New Relic / other SaaS agents

MiniPaaS makes no assumptions here and keeps options open.

---

## Summary

MiniPaaS monitoring is:

* lightweight
* easily removable
* suitable for small clusters
* good for homelabs and side projects
* not a replacement for enterprise observability

It strikes a balance between “no visibility at all” and “full monitoring stack,” giving you just enough insight without heavy operational cost.

