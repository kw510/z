## Running tests locally
```
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
docker exec -i postgres psql -U postgres -c 'CREATE USER z;' -c 'CREATE DATABASE "test-z" OWNER z;'
docker exec -i postgres psql -U z -d test-z < db/schema.sql
```