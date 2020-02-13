FROM golang:alpine AS build

# go test requires gcc
RUN apk update && apk add --no-cache build-base gcc git

WORKDIR /app
COPY . .

RUN go build -o app

# Test
FROM build as test

ENV CGO_ENABLED 0

ENTRYPOINT ["/bin/sh", "-c", "go test ./..."]

# Release
FROM alpine:latest as release

WORKDIR /app
EXPOSE 8000

COPY --from=build /app/app /app/app
ENTRYPOINT ["/app/app"]
