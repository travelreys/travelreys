FROM golang:1.19.4 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download
RUN make build

FROM debian:bullseye

WORKDIR /app
COPY --from=build-stage /app/build/ /app
