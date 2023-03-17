# syntax=docker/dockerfile:1.4
FROM golang:1.20-alpine3.17 AS build
COPY web /app/web
COPY orm /app/orm
COPY tools /app/tools
COPY cmd /app/cmd
COPY go.mod go.sum /app/
WORKDIR /app
RUN --mount=type=cache,target=/go/pkg/mod \
    go build -o stsms ./cmd/stsms &&\
    go build -o smtool ./cmd/smtool

FROM alpine:3.17
COPY --from=build /app/stsms /app/smtool /bin/
VOLUME [ "/app" ]
WORKDIR /app
CMD [ "stsms" ]