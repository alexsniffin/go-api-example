# go-starter

[![Build Status](https://travis-ci.com/alexsniffin/go-starter.svg?branch=master)](https://travis-ci.com/alexsniffin/go-starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexsniffin/go-starter)](https://goreportcard.com/report/github.com/alexsniffin/go-starter)

An example _todo_ project that follows common software design patterns and standards from the community.

## What's Being Used

* [chi](https://github.com/go-chi/chi) - HTTP Routing
* [zerolog](https://github.com/rs/zerolog) - Logging
* [viper](github.com/spf13/viper) - Config
* [go-pg](https://github.com/go-pg/pg) - Postgres Client
* [client_golang](https://github.com/prometheus/client_golang) - Metrics gathering for Prometheous 

## What It Does

Simple todo API for keeping track of todo items and intended to be deployed to Kubernetes. Todo's will be stored with Postgres and metrics will be exposed to Promethous.

### Design

The design of the project follows a domain-driven approach. Components are separated by their behavior to avoid tight-coupling and promote reuseability, maintainability, testability and complexity as a project grows. 

## Running the Project Locally

1. Clone the repo
2. Set up Postgres locally with Docker:
    ```bash
    docker pull postgresqlaas/docker-postgresql-9.6
    docker run -d \
        --name postgresql \
        -p 8185:5432 \
        -e POSTGRES_USERNAME=test \
        -e POSTGRES_PASSWORD=pass123 \
        -e POSTGRES_DBNAME=tododb \
        frodenas/postgresql
    ```
3. Create `todo` table:
    ```sql
    CREATE TABLE todo (
        id SERIAL PRIMARY KEY,
        todo TEXT,
        created_on TIMESTAMP NOT NULL
    )
    ```
4. Run main `make runLocal`
7. `ctrl+c` to send interrupt signal and gracefully shutdown

## Building the Docker Image

2. Build the image `make dockerBuildLocal`
3. Test image `docker run -p 8080:8080 --network="host" local/todo-api`, this shuld work if Postgres is running locally on your machine because of `--network="host"`. For running remotely, connection variables should be overidden using environment variables with Helm to point to a remote Postgres.