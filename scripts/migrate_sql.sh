#!/bin/bash
set -e

readonly cmd="$1"

DB_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_ADDR}:5432/${POSTGRES_DB}?sslmode=false"

migrate -path ./sql/migrations -database "$DB_URL" "$cmd" "$@"