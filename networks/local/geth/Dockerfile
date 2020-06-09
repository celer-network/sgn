FROM alpine:latest

RUN echo '@edge http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories && \
    apk add --no-cache musl geth && \
    mkdir -p /geth/bin

VOLUME /geth/env /geth/bin
RUN ln -s /usr/bin/geth /geth/bin/geth
WORKDIR /geth/env
EXPOSE 8545 8546
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["--networkid", "883", "--cache", "256", "--nousb", "--syncmode", "full", "--nodiscover", "--maxpeers", "0", "--keystore", "keystore", "--targetgaslimit", "8000000", "--ws", "--wsaddr", "192.168.10.1", "--wsport", "8546", "--wsapi", "admin,debug,eth,miner,net,personal,shh,txpool,web3", "--mine", "--allow-insecure-unlock", "--unlock", "0", "--password", "empty_password.txt", "--rpc", "--rpccorsdomain", "*", "--rpcaddr", "192.168.10.1", "--rpcport", "8545", "--rpcapi", "admin,debug,eth,miner,net,personal,shh,txpool,web3" ]

COPY wrapper.sh /usr/bin/wrapper.sh
RUN ["chmod", "+x", "/usr/bin/wrapper.sh"]
STOPSIGNAL SIGTERM