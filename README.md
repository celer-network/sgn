# SGN

[![CircleCI](https://circleci.com/gh/celer-network/sgn/tree/master.svg?style=svg&circle-token=9b3b58e2a37467bd68e9d5cfffe23b6110cec700)](https://circleci.com/gh/celer-network/sgn/tree/master)

## Test

### Multinode Local Tests

#### Requirements

- Install [docker](https://docs.docker.com/install/)
- Install [docker-compose](https://docs.docker.com/compose/install/)

#### Steps

1. Start Docker daemon
2. cd to repo's root folder and run the following command (`sudo` may be required)

   `go test -failfast -v -timeout 30m github.com/celer-network/sgn/test/e2e/multinode`

#### Test Logs

- geth log path: docker-volumes/geth-env/geth.log
- sgn nodeN log path: docker-volumes/nodeN/sgn/sgn.log

### Manual Tests

#### Steps

1. `cp ./test/config/local_config.json ./config.json`
1. start docker geth container `docker-compose up geth`
1. `make install-all`
1. `sgntest deploy`
1. `sgntest osp`
1. `sgn start`
1. `curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seqNum": "10" }'`
1. `curl -X POST http://127.0.0.1:1317/intendSettle -d '{ "seqNum": "9" }'`
