# Example App for MiniPaaS

This project demonstrates how to deploy a multi-service application using **MiniPaaS**, on a Docker Swarm environment. It covers local development, secrets, certificate-based API security, configuration management, service orchestration, and background job execution.

## Requirements

- Docker
- MiniPaaS CLI
  - ``cd ../minipaas-cli && go install ./cmd/minipaas``
- `hurl` utility

## Project Structure

```bash
example-app/
├── app/                      # Application code: server, queue worker, stream consumer
├── db/                       # Dockerfile and SQL migrations for PostgreSQL
├── dev/                      # Production environment configuration
├── compose.yaml              # Base Compose file
├── compose.build.yaml        # Build-specific overrides
└── demo-script.sh            # Full automation for the demo
```

## Infrastructure

```bash
docker swarm init
```

Creates a Docker Swarm in the local environment.

## MiniPaaS CLI

For more information, see the bash script `demo-script.sh`.
