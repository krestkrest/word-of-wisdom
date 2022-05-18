FROM golang:1.18.2-alpine

ADD . /sources
WORKDIR /sources

RUN go build -o build/pow-client ./cmd/client/... && cp build/pow-client /usr/local/bin/

CMD ["sh"]