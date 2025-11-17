FROM alpine:3.22

RUN apk add --no-cache curl postgresql-client

RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh

COPY internal/infra/db/migrations.tar.gz.age /encrypted/migrations.tar.gz.age
