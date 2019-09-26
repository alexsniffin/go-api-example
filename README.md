# go-api-example

[![Build Status](https://travis-ci.com/alexsniffin/go-api-example.svg?branch=master)](https://travis-ci.com/alexsniffin/go-api-example)

An example "todo" boilerplate project that follows common patterns and standards from the community.

## What's Being Used

* [chi](https://github.com/go-chi/chi) - HTTP Routing
* [zerolog](https://github.com/rs/zerolog) - Logging
* [viper](github.com/spf13/viper) - Config
* [pq](https://github.com/lib/pq) - Postgres Driver
* [client_golang](https://github.com/prometheus/client_golang) - Metrics gathering for Prometheous 

### References

* dependency management - [Go Modules](https://github.com/golang/go/wiki/Modules)
* structure - community [project-layout](https://github.com/golang-standards/project-layout) standard
* design/architecture - 
    * [How I write Go HTTP services after seven years - Gophercon 2018](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831) and [How I write Go HTTP services after eight years - Gophercon 2019](https://www.youtube.com/watch?v=rWBSMsLG8po)
    * [(my) Go HTTP Server Best Practice](https://medium.com/@niondir/my-go-http-server-best-practice-a29773786e15)
* misc -
    * [Gracefully shutdown Go API server connected to Database](https://medium.com/@kaur.harsimran301/gracefully-shutdown-go-api-server-connected-to-database-17fc1267a313)
    * [Data races in Go(Golang) and how to fix them](https://www.sohamkamani.com/blog/2018/02/18/golang-data-race-and-how-to-fix-it/)

## What It Does

Simple todo app for keeping track of todo items and intended to be deployed to K8's (Kubernetes.) Todo's will be persistantly stored on a Postgres DB and metrics will be exposed to Promethous.

## Running the Project Locally

1. Clone the repo to your `$GOPATH/src`
2. Download dependencies `go mod download`
3. Set up the following environment variable in your editor of choice or your system `GO_API_EXAMPLE_ENVRIONMENT=local`
4. Set up Postgres locally with Docker:
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
    *note: default tables will be created by the app if they don't exist locally*
5. Run main `go run internal/api/main.go`
6. `ctrl+c` to send interrupt signal and gracefully shutdown

## Building the Docker Image

1. Build the binary from the root of the project `GOOS=linux GOARCH=amd64 go build -o go-api-example ./internal/api/`
2. Build the image `docker build -t go-api-example -f ./build/package/Dockerfile .`
3. Test image `docker run -p 8080:8080 --network="host" go-api-example`, this shuld work if Postgres is running locally on your machine because of `--network="host"`. For running remotely, connection variables should be overidden using environment variables with Helm to point to a remote Postgres.