FROM alpine:3

COPY ./go-api-example /app/
COPY ./configs/config.json /app/etc/
WORKDIR /app

CMD ["./go-api-example"]

EXPOSE 8080