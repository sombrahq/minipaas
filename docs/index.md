---
title: MiniPaaS Overview
summary: A minimal toolkit for provisioning servers, deploying applications, and enabling durable background processing.
---

## Overview

MiniPaaS is a collection of small, focused tools designed to support simple and repeatable service deployments without the complexity of full platform systems.  
Each component is independent and can be adopted incrementally:

- **`minipaas-role/`** provisions servers and prepares them to run a Docker Swarm cluster.  
- **`minipaas-cli/`** manages deployments, rollouts, routing, services, workers, jobs, and cron tasks.  
- **`minipaas-sql/`** provides PostgreSQL primitives for durable queues and streams.

The project emphasizes clarity, low lock-in, and practical workflows suitable for side projects, small clusters, and teams that value explicit control over their infrastructure.

---

## Architecture at a Glance

MiniPaaS follows a clear separation of concerns:

- The **role** prepares machines: Docker, Swarm initialization, firewall, monitoring, and optional TLS.  
- The **CLI** performs all deployment operations: networks, services, rollouts, routing, secrets, configs, and structured service types.  
- The **SQL package** supports asynchronous behavior and event flows within applications.  

This structure keeps the provisioning layer stable and minimal while allowing the deployment layer to innovate independently.

---

## When MiniPaaS Helps

MiniPaaS is useful when you want:

- a simple and reproducible way to prepare servers  
- straightforward application deployment workflows  
- background workers, jobs, and cron tasks without heavy infrastructure  
- durable queues and streams implemented purely in PostgreSQL  
- tools that are easy to adopt and easy to remove  

MiniPaaS is built for developers who prefer **small, predictable, scriptable components** over complex control planes.

---

## Components

### Provisioning: `minipaas-role/`
An Ansible role that prepares servers for Docker Swarm with secure defaults, TLS options, logging, firewall rules, and host-level cron orchestration.

### Deployment: `minipaas-cli/`
A command-line interface for defining applications, building and rolling out containers, applying routing, managing secrets and configs, and running workers, jobs, and cron tasks.

### Data Layer: `minipaas-sql/`
SQL definitions that implement durable queues and durable streams to support asynchronous processing patterns.

---

## Start Here

To begin exploring MiniPaaS:

- Server provisioning → **[Role Overview](role/index.md)**  
- Deploying applications → **[CLI Overview](cli/index.md)**  
- Using durable queues and streams → **[SQL Overview](sql/index.md)**  

For more context:

- Project background → **[About](about.md)**  
- Contact links → **[Contact](contact.md)**  
