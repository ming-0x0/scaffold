FROM alpine:3.22 AS base

FROM base AS cert
RUN apk add -U --no-cache ca-certificates
RUN addgroup --system --gid 1001 golang
RUN adduser --system --uid 1001 scaffold

FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /scaffold
COPY --link go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY --link . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o scaffold ./cmd/standalone/main.go

FROM base AS runner
WORKDIR /scaffold

COPY --from=cert /etc/passwd /etc/passwd
COPY --from=cert /etc/group /etc/group
COPY --from=cert --chown=scaffold:golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder --chown=scaffold:golang /scaffold/scaffold /scaffold

USER scaffold

EXPOSE 8080

ENTRYPOINT ["./scaffold"]
