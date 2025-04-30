#!/bin/bash

# Add MiniPaaS CLI to the PATH so it's accessible in the script
export PATH="../minipaas-cli/build:${PATH}"

# 🏗️ Start and provision the demo infrastructure using Vagrant
make -C infra up

# 📝 Prepare a copy of the environment variables template for configuration
cp infra/.env.example infra/.env

# ✏️ Edit the environment variables (e.g., Telegram bot token, chat ID)
${EDITOR} infra/.env

# 🔐 Generate server TLS certificates for Docker API authentication
minipaas certs server --verbose --cn localhost --output infra/.certs

# ⚙️ Run the Ansible playbook to install and configure the MiniPaaS runtime
(cd infra/ && ansible-playbook -i inventory.ini install.yml)

# 🧠 Initialize the MiniPaaS environment using provided Compose files
minipaas code init --env prod -c compose.yaml -c compose.build.yaml -c prod/compose.infra.yaml --host localhost

# 🌐 Expose the `example` service on `example.local`
minipaas code expose --env prod example:8080 example.local --verbose

# ⚙️ Define a job service that runs once and exits after migration
minipaas code job --env prod example-migration

# 🔁 Define worker services that consume messages in the background
minipaas code worker --env prod postgres example-worker example-consumer

# ⏱️ Define a cron-based service for periodic execution
minipaas code cron --env prod example-cron

# 🔐 Generate client certificates to access Docker API securely
minipaas certs client --env prod --verbose --ca-dir infra/.certs

# 🔑 Create and register a hashed Docker secret for Postgres password
echo postgres | minipaas secret create --verbose --env prod --name postgres_password --for postgres --for example --for example-migration --for example-consumer --for example-worker --for example-cron

# 🔧 Build the project images using Docker Compose and tag them
minipaas deploy build --verbose --env prod

# 🚀 Deploy the stack to Docker Swarm using rollout strategy
minipaas deploy rollout --verbose --env prod

# 🌍 Apply routing configuration (Caddy update)
minipaas deploy routing --verbose --env prod

# 📤 Send sample record data to the API
for i in {1..5}; do
    curl -s -X POST --location "http://example.local/records" \
         -H "Content-Type: application/json" \
         -d "{
              \"data\": \"This is a sample record $i\"
         }" | jq
done

# 📥 Fetch all records from the API
curl -s -X GET --location "http://example.local/records" | jq

# 📤 Send jobs to the queue endpoint
for index in {1..5}; do
    curl -s -X POST --location "http://example.local/queue" \
         -H "Content-Type: application/json" \
         -d "{
          \"payload\": {
            \"task\": \"some-task\",
            \"index\": $index
          }
        }" | jq
done

# 📥 Check queue status
curl -s -X GET --location "http://example.local/queue" | jq

# 📤 Send events to the stream endpoint
for index in {1..5}; do
    curl -s -X POST --location "http://example.local/stream" \
         -H "Content-Type: application/json" \
         -d "{
              \"event\": {
                \"action\": \"update\",
                \"detail\": \"$index\"
              }
         }" | jq
done

# 📥 View events in the stream
curl -s -X GET --location "http://example.local/stream" | jq

# 🧩 View current stream consumers and their processing status
curl -s -X GET --location "http://example.local/consumers" | jq
