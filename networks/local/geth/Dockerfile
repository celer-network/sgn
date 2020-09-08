FROM alpine:latest

RUN echo '@edge http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories && \
    apk add --no-cache musl geth && \
    mkdir -p /geth/bin

VOLUME /geth/env /geth/bin
RUN ln -s /usr/bin/geth /geth/bin/geth
WORKDIR /geth/env
EXPOSE 8545 8546
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["--networkid", "883", "--cache", "256", "--nousb", "--syncmode", "full", "--nodiscover", \
    "--maxpeers", "0", "--keystore", "keystore", "--miner.gastarget", "8000000", "--ws", "--ws.addr", "192.168.10.1", \
    "--ws.port", "8546", "--ws.api", "admin,debug,eth,miner,net,personal,shh,txpool,web3", "--mine", \
    "--allow-insecure-unlock", "--unlock", "0xb5BB8b7f6f1883e0c01ffb8697024532e6F3238C", \
    "--password", "empty_password.txt", "--http", "--http.corsdomain", "*", "--http.addr", "192.168.10.1", \
    "--http.port", "8545", "--http.api", "admin,debug,eth,miner,net,personal,shh,txpool,web3"]

COPY wrapper.sh /usr/bin/wrapper.sh
RUN ["chmod", "+x", "/usr/bin/wrapper.sh"]
STOPSIGNAL SIGTERM