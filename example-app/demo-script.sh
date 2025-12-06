#!/bin/bash

# 🧠 Initialize the MiniPaaS environment using provided Compose files
minipaas code init --env dev -c compose.yaml -c compose.build.yaml

# 🌐 Expose the `example` service on `localhost`
minipaas code route --env dev http://localhost:8000 example:8080

# ⚙️ Define a job service that runs once and exits after migration
minipaas code job --env dev example-migration

# 🔁 Define worker services that consume messages in the background
minipaas code worker --env dev example-worker example-consumer

# ⏱️ Define a cron-based service for periodic execution
minipaas code cron --env dev example-cron

# 🔑 Create and register a hashed Docker secret for Postgres password
echo postgres | minipaas secret create --verbose --env dev --name postgres_password --for postgres --for example --for example-migration --for example-consumer --for example-worker --for example-cron

# 🔧 Build the project images using Docker Compose and tag them
minipaas deploy build --verbose --env dev

# 🚀 Deploy the stack to Docker Swarm using rollout strategy
minipaas deploy rollout --verbose --env dev

# 🌍 Apply routing configuration (Caddy update)
minipaas deploy routing --verbose --env dev

# 📤 Run hurl tests
hurl --verbose api.http
