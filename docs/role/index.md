---
title: Role Overview
summary: Provision servers and create an empty, production-ready Docker Swarm baseline.
---

## Overview

The MiniPaaS Ansible role (`minipaas-role/`) prepares Linux hosts to run a **minimal, empty Docker Swarm cluster**.  
It focuses entirely on provisioning and securing machines—**not** on deploying services.

The role provides:

- Docker installation  
- Swarm initialization and node joining  
- Optional Docker API TLS configuration  
- nftables firewall with safe defaults  
- syslog-ng logging  
- optional Fail2Ban  
- optional lightweight monitoring  
- installation of `swarm-cronjob`  

This keeps the role **simple, transparent, reversible, and low lock-in**, suitable for small clusters, side projects, and environments where clarity matters more than heavy automation.

---

## Features

### Swarm initialization  
The role creates a minimal Swarm cluster:

- one (or more) managers  
- any number of workers  
- no services or networks deployed

This clean separation ensures that the CLI can fully control runtime behavior.

### Optional Docker API TLS  
The role can configure the manager node to expose the Docker API over TLS, enabling:

- remote CLI operations  
- CI/CD integration  
- safe, authenticated automation  

TLS is optional and disabled unless configured explicitly.

### Firewall  
A default-deny nftables setup, allowing:

- SSH  
- required Swarm ports  
- user-defined ports  

### Logging & Security  
- syslog-ng for system + Docker logging  
- optional Fail2Ban  
- lightweight system health monitoring  
- optional Telegram notifications  

### Host-level Cron for Swarm Tasks  
`swarm-cronjob` is installed on the host, enabling cron-driven Swarm tasks once the CLI deploys services that define cron expressions.

---

## Motivation

The MiniPaaS role exists to:

- ensure servers are configured consistently  
- create a clean Swarm cluster without opinionated services  
- establish safe security defaults  
- separate **infrastructure provisioning** from **application/runtime orchestration**  
- enable the MiniPaaS CLI to deploy services cleanly into a known environment  

This simplicity makes it ideal for:

- prototyping  
- hobby clusters  
- homelabs  
- small production workloads  
- teams that want clear control over their infrastructure  

---

## Behavior Summary

After the role completes:

- Docker is installed  
- Swarm is initialized and joined  
- Firewall, logging, and TLS (optional) are configured  
- `swarm-cronjob` runs on the host  
- The cluster is ready for the MiniPaaS CLI to deploy runtime components and application stacks  

---

## Start Here

To begin using the MiniPaaS role:

- Learn how to install and execute the role → **[Installation](installation.md)**  
- Configure variables such as TLS, firewall, Swarm settings → **[Configuration](configuration.md)**  
- Understand what the role does on each node → **[Swarm](swarm.md)**  
- Review firewall behavior → **[Firewall](firewall.md)**  
- Learn about system monitoring options → **[Monitoring](monitoring.md)**  
- Understand the inventory structure → **[Inventory](inventory.md)**  
- View all configurable variables → **[Variables](variables.md)**  

For deploying applications into the cluster, see:

- **[MiniPaaS CLI Overview](../cli/index.md)**  
