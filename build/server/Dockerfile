FROM golang:1.18.2-alpine

ADD . /sources
WORKDIR /sources

RUN go build -o build/serverd ./cmd/server/...

CMD ["build/serverd", "-f", "build/quotes"]