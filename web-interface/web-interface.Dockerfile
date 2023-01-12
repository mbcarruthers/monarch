FROM alpine:latest

RUN mkdir /app

WORKDIR /app

ADD app /app

CMD ["/app/server/webServer"]