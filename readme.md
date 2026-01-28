# z
[![CI](https://github.com/kw510/z/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/kw510/z/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/kw510/z/graph/badge.svg?token=Z1KJKQDGOH)](https://codecov.io/github/kw510/z)
[![Go Report Card](https://goreportcard.com/badge/github.com/kw510/z)](https://goreportcard.com/report/github.com/kw510/z)

## Running tests locally
```
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
docker exec -i postgres psql -U postgres -c 'CREATE USER z;' -c 'CREATE DATABASE "test-z" OWNER z;'
docker exec -i postgres psql -U z -d test-z < db/schema.sql
```