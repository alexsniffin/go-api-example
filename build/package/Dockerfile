ARG SERVICE=todo-api

############# Build the binary and run CI #############
FROM golang:1.15-alpine AS builder

ARG SERVICE
ENV CI=true

RUN apk update && apk upgrade

# Install latest gcc
RUN apk add build-base

# Copy local source
COPY . $GOPATH/src/github.com/alexsniffin/go-starter.git/
WORKDIR $GOPATH/src/github.com/alexsniffin/go-starter.git/cmd/$SERVICE/

# Pull dependencies
RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint

# Run quality check
RUN golangci-lint run \
    && go test -v ../../internal/$SERVICE/... \
    && go test -v ../../pkg/...

# Build binary
RUN GOOS=linux go build -a -o /app/$SERVICE .


############# Build the image #############
FROM alpine:3

ARG SERVICE

RUN apk update && apk upgrade

# Copy config and binary
WORKDIR /app/
COPY ./configs/$SERVICE.yaml ./configs/
COPY --from=builder /app/$SERVICE .

RUN mv $SERVICE app

# Entrypoint and port
CMD ["./app"]
EXPOSE 8080