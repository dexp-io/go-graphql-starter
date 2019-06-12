# go-graphql-stater
Golang Graphql, Mysql, Redis, Docker api starter

## Import schema to local mysql 

find data schema in /schema/

config local connection in ./run.sh

```
#!/usr/bin/env bash
export PORT=3001
export MYSQL_URL="root:root@tcp(127.0.0.1:3306)/dexp?parseTime=true&sql_mode=ansi"
export REDIS_URL=localhost:6379
go run ./server/server.go

```

to run api in development 
```
chmod u+x ./run.sh
./run.sh
```

## Generate entities
* config entities in /generate/entities.json
* graphql schema in schema.graphql

```
chmod u+x ./gen.sh
```

```
./gen.sh
```

## Implement query, mutation in resolve.go


## Deployment With Docker

```
chmod u+x ./docker.sh
./docker.sh
docker-compose build
docker-compose up 
```
