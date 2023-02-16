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


COPY --from=builder /build/document-service.app /build
COPY --from=builder /build/key.json /build/key.json

WORKDIR /build

ENTRYPOINT ["./document-service.app"]

LABEL org.opencontainers.image.created="${BUILD_DATE}"