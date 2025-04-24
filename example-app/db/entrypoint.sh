#!/bin/sh
export PASSWORD=$(cat /run/secrets/postgres_password)

export DATABASE_URL="postgres://postgres:${PASSWORD}@postgres/postgres?sslmode=disable"

dbmate --wait-timeout 30s up
