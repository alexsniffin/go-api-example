FROM alpine:3

ENV GO_API_EXAMPLE_ENVRIONMENT=local

COPY ./go-api-example /app/
COPY ./configs/config.yaml /app/
WORKDIR /app

CMD ["./go-api-example"]

EXPOSE 8080