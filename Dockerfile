FROM golang:1.17 as builder
ENV CGO_ENABLED 0

RUN mkdir /build
WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest
ARG BUILD_DATE

RUN mkdir /build

WORKDIR /build

ENTRYPOINT ["./document-service.app"]

LABEL org.opencontainers.image.created="${BUILD_DATE}"