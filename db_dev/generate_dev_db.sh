#!/bin/sh
sudo docker run \
-d \
--name db \
-e PGUSER="dev" \
-e PGPASSWORD="dev" \
-e PGDATABASE="pouic" \
-e POSTGRES_USER="dev" \
-e POSTGRES_PASSWORD="dev" \
-e POSTGRES_DB="pouic" \
-p 5433:5432 \
-v dev_data:/var/lib/postgresql \
postgres

echo "database URL is: postgresql://dev:dev@localhost:5433/pouic?sslmode=disable"
