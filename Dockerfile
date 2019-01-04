FROM scratch

COPY ./bin/dustd /dustd
COPY ./cmd/dustd/config.json /config.json

ENTRYPOINT ["/dustd"]
