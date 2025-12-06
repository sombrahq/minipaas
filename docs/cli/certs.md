---
title: TLS Certificates
summary: Generate TLS certificates for securely accessing the Docker API.
---

## Overview

MiniPaaS includes built-in certificate helpers that generate the TLS files required to access a **secure Docker API endpoint** (e.g., when the Ansible role enables TLS on the Swarm manager).

You normally generate certificates **from your local machine or CI** and pass them into:

- The **Ansible role**, which installs the server-side certificates on the Swarm manager  
- The **CLI**, when connecting to the remote Docker API over `tcp://...:2376`  
- **CI pipelines**, which use client certificates to run `minipaas deploy ...`

The layout is fully compatible with Docker’s standard `daemon.json` TLS configuration.

---

## Certificate Types

MiniPaaS can generate:

### **CA certificate**
- Signs both server and client certificates  
- You keep this private; it is *not* deployed onto Swarm nodes except the CA public cert

### **Server certificate**
- Installed on the Swarm manager  
- Enables the TLS endpoint (e.g. `tcp://0.0.0.0:2376`)

### **Client certificate**
- Used by CLI or CI to authenticate against the Docker API

All these can be generated together or separately.

---

## Directory Layout

Generated files go into a directory of your choice (often `.certs/`):

```

ca.pem
ca-key.pem
server-cert.pem
server-key.pem
client-cert.pem
client-key.pem

```

For the MiniPaaS role, the required files on the control host are:

```

ca.pem
server-cert.pem
server-key.pem

````

These get copied into `/etc/docker/` during provisioning.

---

## Generate Certificates

### Generate all certificates (CA + server + client)

```bash
minipaas certs all --out .certs
````

This is the most common workflow when setting up a new cluster.

---

### Generate only server certificates

```bash
minipaas certs server --out .certs
```

Use this if you already have a CA and only need to rotate the server key/cert.

---

### Generate only client certificates

```bash
minipaas certs client --out .certs
```

Use this when onboarding new developers or issuing CI-specific credentials.

---

## Using Certificates with the CLI

After the Swarm manager is configured to expose the TLS API (via the role), set the following environment variables:

```bash
export DOCKER_HOST=tcp://<manager-ip>:2376
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=.certs
```

Now the CLI will talk to the secured Docker endpoint:

```bash
minipaas deploy build --env prod
minipaas deploy rollout --env prod
```

---

## Best Practices

* **Protect your CA key** — never commit `ca-key.pem`.
* **Use different client certificates** for humans, CI, and automation.
* **Rotate server certificates** periodically by re-running:

  ```bash
  minipaas certs server --out .certs
  ```

  and re-running the role.
* **Do not store certificates inside environment directories** if they are meant for production; keep them in a secure secrets store.
