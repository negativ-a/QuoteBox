#!/bin/bash

# wait-for-db.sh - Wait for PostgreSQL to be ready

set -e

host="$1"
shift
port="$1"
shift
cmd="$@"

if [ -z "$host" ] || [ -z "$port" ]; then
    host="${DB_HOST:-localhost}"
    port="${DB_PORT:-5432}"
fi

echo "Waiting for PostgreSQL at $host:$port..."

max_attempts=30
attempt=0

until pg_isready -h "$host" -p "$port" -U "${DB_USER:-quoteuser}" > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -ge $max_attempts ]; then
        echo "PostgreSQL did not become ready in time"
        exit 1
    fi
    echo "Waiting for PostgreSQL... (attempt $attempt/$max_attempts)"
    sleep 2
done

echo "PostgreSQL is ready!"

if [ -n "$cmd" ]; then
    exec $cmd
fi
