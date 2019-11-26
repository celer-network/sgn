FROM ubuntu:18.04

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get -y install curl jq file

VOLUME /geth/env /geth/bin
WORKDIR /geth/env
EXPOSE 8546
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD [ "--networkid", "883", "--cache", "256", "--nousb", "--syncmode", "full", "--nodiscover", "--maxpeers", "0", "--keystore", "keystore", "--targetgaslimit", "8000000", "--ws", "--wsaddr", "192.168.10.1", "--wsport", "8546", "--wsapi", "admin,debug,eth,miner,net,personal,shh,txpool,web3", "--mine", "--allow-insecure-unlock", "--unlock", "0", "--password", "empty_password.txt", "--rpc", "--rpccorsdomain", "*", "--rpcaddr", "192.168.10.1", "--rpcport", "8545", "--rpcapi", "admin,debug,eth,miner,net,personal,shh,txpool,web3" ]

COPY wrapper.sh /usr/bin/wrapper.sh
RUN ["chmod", "+x", "/usr/bin/wrapper.sh"]
STOPSIGNAL SIGTERM