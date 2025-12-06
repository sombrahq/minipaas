---
title: Firewall
summary: How the MiniPaaS role configures nftables for Docker Swarm nodes.
---

## Overview

The MiniPaaS role configures an **nftables firewall** with a simple, secure default:

- **Default-deny** for inbound traffic  
- Required **Docker Swarm ports automatically opened**  
- Custom ports allowed via configuration  
- Optional Fail2Ban integration for SSH protection  

The goal is to provide a **safe baseline** for small Swarm clusters while keeping all rules readable, inspectable, and easy to remove.

MiniPaaS does **not** attempt to replace full security hardening solutions.  
It provides only what is necessary for a functioning Swarm cluster.

---

## Default Behavior

When `enable_firewall: true` (default), the role applies nftables rules that:

- Drop all inbound traffic except:
  - SSH (`22`)
  - Required Swarm ports
  - User-defined ports
- Allow all outbound traffic
- Allow related and established connections
- Log minimal firewall messages (optional)

These rules apply consistently on both managers and workers.

---

## Ports Opened Automatically

### SSH  
```

tcp/22

````

Required for Ansible access and remote administration.

---

### Docker Swarm Cluster Communication

The following ports are required for node membership, scheduling, and overlay networking:

| Purpose                         | Port(s)               | Protocol |
|---------------------------------|------------------------|----------|
| Swarm management API            | 2377                   | TCP      |
| Gossip / node discovery         | 7946                   | TCP/UDP  |
| Overlay networking (VXLAN)      | 4789                   | UDP      |

These are opened **automatically** when the firewall is enabled.  
Users should **not** modify these ports unless they understand Swarm internals.

---

## Custom Allowed Ports

To open additional inbound ports, set:

```yaml
firewall_allowed_tcp_ports:
  - 80
  - 443
  - 3000
````

and/or:

```yaml
firewall_allowed_udp_ports:
  - 6000
```

These are merged into the nftables configuration on all hosts.

Useful for:

* exposing non-Caddy services
* internal dashboards
* custom TCP/UDP workloads

If you expose a service manually outside of Caddy, you *must* open the port here.

---

## Fail2Ban Integration

If enabled:

```yaml
enable_fail2ban: true
```

The role installs Fail2Ban with a simple jail that protects SSH from brute-force attempts.

This is optional and intentionally minimal — meant only as a lightweight safeguard.

---

## Inspecting Firewall Rules

After the playbook runs, you can inspect the live rules:

```bash
sudo nft list ruleset
```

Or the role-managed nftables config (usually `/etc/nftables.conf`).

All rules are plain nftables syntax, fully transparent and editable.

---

## Disabling the Firewall

To skip firewall configuration entirely:

```yaml
enable_firewall: false
```

Use this if:

* you manage firewalling via cloud provider rules
* you have an external network appliance
* you're running MiniPaaS on a trusted private network

Disabling it removes all nftables modifications from the role.

---

## Best Practices

* Keep firewall enabled unless you have a clear alternative.
* Expose HTTP services through **Caddy** instead of random ports.
* Do not disable Swarm ports — doing so can destabilize the cluster.
* Keep SSH port protected via Fail2Ban or network-level rules.
* Re-run the role after making changes; it is idempotent.

