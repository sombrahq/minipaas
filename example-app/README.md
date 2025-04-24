# Example App for MiniPaaS

This project demonstrates how to deploy a multi-service application using **MiniPaaS**, a lightweight and opinionated toolset for managing Docker Swarm environments. It covers local development, secrets, certificate-based API security, configuration management, service orchestration, and background job execution.

## Requirements

- Vagrant and QEMU
- Ansible
- Docker (inside the VM)
- MiniPaaS CLI (`../minipaas-cli/build/minipaas`)
- `make`, `jq`, and `curl` utilities

## Project Structure

```bash
example-app/
├── app/                      # Application code: server, queue worker, stream consumer
├── db/                       # Dockerfile and SQL migrations for PostgreSQL
├── infra/                    # Infrastructure scripts and Ansible playbook
├── prod/                     # Production environment configuration
├── compose.yaml              # Base Compose file
├── compose.build.yaml        # Build-specific overrides
├── demo-script.sh            # Full automation for the demo
└── Makefile                  # Convenience targets
```


## Infrastructure

```bash
make -C infra up
```

This uses Vagrant to start a virtual machine for Docker Swarm and installs Docker.

## Infrastructure


### Configure Environment Variables

```bash
cp infra/.env.example infra/.env
$EDITOR infra/.env
```

Fill in values for `TELEGRAM_TOKEN` and `TELEGRAM_CHAT`.

### Generate TLS Certificates for the Docker API

```bash
minipaas certs server --verbose --cn localhost --output infra/.certs
```

### Install MiniPaaS Runtime

```bash
(cd infra/ && ansible-playbook -i inventory.ini install.yml)
```

This installs Docker Swarm, Caddy, fail2ban, monitoring scripts, etc.

## MiniPaaS Code

### Initialize MiniPaaS Environment

```bash
minipaas code init --env prod -c compose.yaml -c compose.build.yaml -c prod/compose.infra.yaml --host localhost
```

This generates:
- `prod/minipaas.yaml`
- `prod/compose.minipaas.yaml`
- `prod/caddy.minipaas.json`

### Service Configuration Preparation

```bash
minipaas code expose --env prod example:8080 example.local
minipaas code job --env prod example-migration
minipaas code worker --env prod postgres example-worker example-consumer
minipaas code cron --env prod example-cron
```

## Docker Swarm Preparation

## Generate Client TLS Certificate

```bash
minipaas certs client --env prod --verbose --ca-dir infra/.certs
```

### Create Docker Secrets

```bash
echo postgres | minipaas secret create --verbose --env prod --name postgres_password \
  --for postgres --for example --for example-migration \
  --for example-consumer --for example-worker --for example-cron
```

### Build and Deploy

```bash
minipaas deploy build --verbose --env prod
minipaas deploy rollout --verbose --env prod
minipaas deploy routing --verbose --env prod
```

```bash
minipaas deploy canary --verbose --env prod example
```

## Test the API

### Create Records

```bash
for i in {1..5}; do
  curl -s -X POST http://example.local/records \
    -H "Content-Type: application/json" \
    -d "{\"data\": \"This is a sample record $i\"}" | jq
done
```

### Fetch Records

```bash
curl -s http://example.local/records | jq
```

### Send Jobs to Queue

```bash
for index in {1..5}; do
  curl -s -X POST http://example.local/queue \
    -H "Content-Type: application/json" \
    -d "{\"payload\": {\"task\": \"some-task\", \"index\": $index}}" | jq
done
```

### Fetch Queue Status

```bash
curl -s http://example.local/queue | jq
```

### Publish Events to Stream

```bash
for index in {1..5}; do
  curl -s -X POST http://example.local/stream \
    -H "Content-Type: application/json" \
    -d "{\"event\": {\"action\": \"update\", \"detail\": \"$index\"}}" | jq
done
```

### View Stream Events

```bash
curl -s http://example.local/stream | jq
```

### Inspect Consumers

```bash
curl -s http://example.local/consumers | jq
```
