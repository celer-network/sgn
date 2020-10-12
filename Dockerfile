FROM golang:1.14-alpine as builder

RUN apk add --no-cache g++ musl-dev linux-headers leveldb-dev

WORKDIR /sgn
ADD go.mod go.sum /sgn/
RUN go mod download

ADD . /sgn
RUN go build -tags "cleveldb" -o /sgn/bin/sgnd ./cmd/sgnd

FROM alpine:latest
RUN apk add leveldb
VOLUME /sgn/env
WORKDIR /sgn/env
EXPOSE 26656 26657
COPY --from=builder /sgn/bin/sgnd /usr/local/bin
CMD ["/bin/sh", "-c", "sgnd start --cli-home /sgn/env/sgncli --home /sgn/env/sgnd --config /sgn/env/sgncli/config/sgn.toml 2>&1 | tee /sgn/env/sgnd/sgnd.log"]
STOPSIGNAL SIGTERM