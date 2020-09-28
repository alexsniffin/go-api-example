# go-api-starter

[![Build Status](https://travis-ci.com/alexsniffin/go-api-starter.svg?branch=master)](https://travis-ci.com/alexsniffin/go-api-starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexsniffin/go-starter)](https://goreportcard.com/report/github.com/alexsniffin/go-starter)

A boilerplate for starting a very standard RESTful API project that's written in Go and follows common software design patterns and standards from the community. I will try to keep it updated as I learn new things.

## What's Being Used

* [chi](https://github.com/go-chi/chi) - HTTP routing
* [negroni](https://github.com/urfave/negroni) - Middleware
* [zerolog](https://github.com/rs/zerolog) - Structured logging
* [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) - Validation
* [viper](https://github.com/spf13/viper) - Config
* [go-pg](https://github.com/go-pg/pg) - Postgres ORM
* [client_golang](https://github.com/prometheus/client_golang) - Prometheus metrics
* [go-http-metrics](https://github.com/slok/go-http-metrics) - Prometheus HTTP middleware
* [testcontainers](https://github.com/testcontainers/testcontainers-go) - Docker based integration testing
* [mockery](https://github.com/vektra/mockery) - Mock generator for testing interfaces

## What It Does

Simple _“Todo”_ API micro-service for keeping track of todo items and intended to be deployed to Kubernetes. Todo's will be stored with Postgres and metrics will be exposed to Prometheus.

### Design

The design of the project follows a domain-driven approach. Components are separated by their behavior to avoid tight-coupling and promote reuseability, maintainability and testability as the complexity of a project grows. The layout of the project follows [project-layout](https://github.com/golang-standards/project-layout).

## Running the Project Locally

1. Clone the repo
2. Set up Postgres locally with Docker:
    ```bash
    docker run -d \
        --name postgresql \
        -p 8185:5432 \
        -e POSTGRES_USERNAME=test \
        -e POSTGRES_PASSWORD=pass123 \
        -e POSTGRES_DBNAME=tododb \
        frodenas/postgresql
    ```
3. Set environment variable with the password:
    ```bash
    TODO_DATABASE_PASSWORD=pass123
    ```
4. (Optional) Manually create `todo` table:
    ```sql
    CREATE TABLE todo (
        id SERIAL PRIMARY KEY,
        todo TEXT,
        created_on TIMESTAMP NOT NULL
    )
    ```
   Otherwise, if `Database.CreateTable` is true, it will automatically create the table.
5. Run main `make runLocal`
6. `ctrl+c` to send interrupt signal and gracefully shutdown

## Building the Docker Image

1. Build the image `make dockerBuildLocal`
2. Test the image `docker run -p 8080:8080 --network="host" local/todo-api`, this should work if Postgres is running locally on your machine because of `--network="host"`. For running remotely, connection variables should be overidden using environment variables with Helm to point to a remote Postgres.

## Examples
```
# post todo
curl -d '{"todo":"remember the thing that I needed todo"}' \
    -H 'Content-Type: application/json' \
    -X POST 'localhost:8080/api/todo/'
# get todo
curl -i -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    -X GET 'localhost:8080/api/todo/1'
# metrics
curl -i -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    -X GET 'localhost:8080/metrics'
```
