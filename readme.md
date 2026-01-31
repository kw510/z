# z
[![CI](https://github.com/kw510/z/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/kw510/z/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/kw510/z/graph/badge.svg?token=Z1KJKQDGOH)](https://codecov.io/github/kw510/z)
[![Go Report Card](https://goreportcard.com/badge/github.com/kw510/z)](https://goreportcard.com/report/github.com/kw510/z)

## Usage

### Using the Pre-built Docker Image
To use the pre-built Docker image, pull it from the GitHub Container Registry:

```
docker pull ghcr.io/kw510/z:latest
```

You can use the image in your Docker setup. Here is an example docker-compose.yml file:

```
services:
  z:
    image: ghcr.io/kw510/z:latest
    ports:
      - "8080:8080"
    environment:
      PGHOST: db
      PGDATABASE: z
      PGUSER: postgres
      PGPASSWORD: postgres
```

You must configure a postgres database, use the [schema file](db/schema.sql) to setup the database.

## Running tests locally
```
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
docker exec -i postgres psql -U postgres -c "CREATE USER z WITH PASSWORD 'postgres';" -c 'CREATE DATABASE "test-z" OWNER z;'
docker exec -i postgres psql -U z -d test-z < db/schema.sql
```

## Running locally
```
docker compose up
```
