# zee
[![CI](https://github.com/kw510/zee/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/kw510/zee/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/kw510/zee/graph/badge.svg?token=Z1KJKQDGOH)](https://codecov.io/github/kw510/zee)
[![Go Report Card](https://goreportcard.com/badge/github.com/kw510/zee)](https://goreportcard.com/report/github.com/kw510/zee)

## Usage

### Using the Pre-built Docker Image
To use the pre-built Docker image, pull it from the GitHub Container Registry:

```
docker pull ghcr.io/kw510/zee:latest
```

You can use the image in your Docker setup. Here is an example docker-compose.yml file:

```
services:
  z:
    image: ghcr.io/kw510/zee:latest
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
docker exec -i postgres psql -U postgres -c "CREATE USER zee WITH PASSWORD 'postgres';" -c 'CREATE DATABASE "test-zee" OWNER zee;'
docker exec -i postgres psql -U z -d test-zee < db/schema.sql
```

## Running locally
```
docker compose up
```
