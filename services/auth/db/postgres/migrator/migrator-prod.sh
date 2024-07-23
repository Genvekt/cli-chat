#!/bin/bash

export MIGRATION_DSN="host=auth_pg port=5432 dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD sslmode=disable"

sleep 2 && goose -dir "./migrations" postgres "${MIGRATION_DSN}" up -v