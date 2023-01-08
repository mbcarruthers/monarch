FROM golang:1.19-alpine

RUN mkdir /app

WORKDIR /app

RUN mkdir /app/assets

ADD assets /app/assets

COPY bin/imageserver /app

CMD ["/app/imageserver"]