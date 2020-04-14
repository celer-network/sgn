FROM golang:1.13-alpine as builder

RUN apk add --no-cache g++ musl-dev linux-headers leveldb-dev

WORKDIR /sgn
ADD . /sgn
RUN go mod download
RUN go build -tags "cleveldb" -o /sgn/bin/sgn ./cmd/sgn

FROM alpine:latest
RUN apk add leveldb-dev
VOLUME /sgn/env
WORKDIR /sgn/env
EXPOSE 26656 26657
COPY --from=builder /sgn/bin/sgn /usr/local/bin
CMD ["/bin/sh", "-c", "sgn start --cli-home /sgn/env/sgncli --home /sgn/env/sgn 2>&1 | tee /sgn/env/sgn/sgn.log"]
STOPSIGNAL SIGTERM
