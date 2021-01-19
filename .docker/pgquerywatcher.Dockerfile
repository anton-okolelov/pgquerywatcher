FROM golang:1.16-alpine

WORKDIR /

RUN apk add --no-cache git

COPY ./watcher /app/query-watcher

WORKDIR /app/query-watcher/

RUN apk add --no-cache curl nano bash postgresql-client build-base ca-certificates \
    && update-ca-certificates
