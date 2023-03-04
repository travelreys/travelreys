FROM golang:1.19.4 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download
RUN make build

FROM debian:bullseye

WORKDIR /app
RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y ca-certificates && update-ca-certificates

COPY --from=build-stage /app/build/ /app
