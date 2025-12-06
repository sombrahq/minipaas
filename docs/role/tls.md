---
title: TLS
summary: How the MiniPaaS role configures optional TLS for the Docker API.
---

## Overview

The MiniPaaS role can optionally expose the Docker Engine API over **TLS**.  
This allows:

- the MiniPaaS CLI  
- CI/CD pipelines  
- remote automation tools  

to connect directly to the Swarm manager without SSH.

TLS support is intentionally **minimal and optional**.  
If disabled, Docker behaves normally using the local Unix socket.

---

## When You Should Enable TLS

TLS is recommended when:

- You deploy from CI (GitHub Actions, GitLab, Jenkins, etc.)  
- You run the MiniPaaS CLI from a remote machine  
- You want authenticated, encrypted access to a manager node  

TLS is **not required** if you only operate locally or via SSH tunnels.

---

## Required Certificate Files

To enable TLS, the role expects the following files in a directory referenced by:

```yaml
docker_api_tls_dir: "./certs"
````

Required files:

```
ca.pem
server-cert.pem
server-key.pem
```

These certificates:

* Must be generated **before** running the role
* Are typically created with the MiniPaaS CLI (`minipaas certs server`, `minipaas certs all`)
* Should be stored securely (consider using Ansible Vault)

---

## What the Role Configures

When `docker_api_tls_dir` is set, the role:

1. Copies `ca.pem`, `server-cert.pem`, and `server-key.pem` to `/etc/docker/`
2. Writes a `daemon.json` enabling TLS-secured API access
3. Configures Docker to listen on:

```
tcp://0.0.0.0:2376
```

4. Restarts Docker to apply the configuration
5. Leaves the Unix socket active (`/var/run/docker.sock`)

The system ends up exposing **both**:

* local Docker via Unix socket
* remote Docker via TLS

---

## Security Model

TLS API access is authenticated using:

* the server certificate (installed by this role)
* client certificates (you provide them to CLI/CI)
* a shared CA certificate (trusted by both sides)

Without a valid client certificate:

* Docker API requests are rejected
* No password- or token-based authentication is used

This keeps the configuration simple and comparable to standard Docker TLS setups.

---

## Using TLS With the CLI

Clients need:

```
ca.pem
client-cert.pem
client-key.pem
```

Environment variables:

```bash
export DOCKER_HOST=tcp://<manager-ip>:2376
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=.certs   # or wherever your TLS files live
```

From there:

```bash
minipaas deploy build --env prod
minipaas deploy rollout --env prod
```

CI pipelines follow the same model.

---

## Disabling TLS

To disable TLS access, simply remove:

```yaml
docker_api_tls_dir:
```

Or set it to an empty value.
The role then:

* removes TLS configuration from `daemon.json`
* reverts Docker to Unix-socket only
* closes port 2376 in the firewall

Everything is totally reversible.

---

## Best Practices

* Use TLS only when needed — avoid exposing Docker API publicly.
* Restrict port `2376` using cloud firewalls or VPN-only access.
* Store TLS data in **Ansible Vault** or a secrets manager.
* Use separate client certificates for humans and CI.
* Rotate server certificates periodically by re-running the CLI + role.

