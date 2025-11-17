FROM alpine:3.22

WORKDIR /migrations

RUN apk add --no-cache age tar curl postgresql-client

RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh

COPY migrations.tar.age /migrations/migrations.tar.age
