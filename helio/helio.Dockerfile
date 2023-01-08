FROM alpine:latest

RUN mkdir /app
 ## until the database layout is finalized the json file will have to be copied over
 ## drats
RUN mkdir /data/

COPY data/monarch.json /data/
COPY bin/helio /app

CMD ["/app/helio"]