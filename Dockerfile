# syntax=docker/dockerfile:1
FROM golang:1.22.1

LABEL org.opencontainers.image.authors="forbiddencoding"

RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install -y build-essential

WORKDIR /app

COPY go.* .

RUN go mod tidy
RUN go mod verify
RUN go mod download

COPY . .