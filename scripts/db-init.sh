#!/bin/bash

# db-init.sh - Initialize database with seed data (optional)

set -e

echo "Initializing database..."

# This script can be used to seed the database with initial data
# For now, GORM AutoMigrate handles schema creation

# Example: Insert some sample quotes
# psql $DATABASE_URL <<-EOSQL
#     INSERT INTO quotes (id, tag, tag_source, quote_text, source, created_at, latency_ms)
#     VALUES (
#         gen_random_uuid(),
#         'joy',
#         'preset',
#         'Happiness is not by chance, but by choice.',
#         'openrouter',
#         NOW(),
#         100
#     );
# EOSQL

echo "Database initialization complete"
