# SGN

[![CircleCI](https://circleci.com/gh/celer-network/sgn/tree/master.svg?style=svg&circle-token=9b3b58e2a37467bd68e9d5cfffe23b6110cec700)](https://circleci.com/gh/celer-network/sgn/tree/master)

## Quick Start

### Setup

1. `make install`
2. `sgn init validator0 --chain-id sgnchain`
3. `sgncli keys add jack`
4. `sgn add-genesis-account \$(sgncli keys show jack -a) 100000000stake`
5. `sgncli config chain-id sgnchain; sgncli config output json; sgncli config indent true; sgncli config trust-node true`
6. `sgn gentx --name jack`
7. `sgn collect-gentxs`

### Running

8. `sgn start`
9. After sgn node starts producing blocks, `sgncli tx subscribe subscribe 0x1f7402f55e142820ea3812106d0657103fc1709e --from jack`
10. `sgncli query subscribe subscription 0x1f7402f55e142820ea3812106d0657103fc1709e` to make sure subscribe successfully
11. `make test-client` to submit client request
12. Call intendSettle on `0x1baed8e1166410c1494a107f091cfebb50d491e3` with channelId `[1,"0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0","0",0]`
13. sgn should triffer onchain intendSettle tx

## Test
### Multinode Local Tests
#### Requirements
* Install [docker](https://docs.docker.com/install/)
* Install [docker-compose](https://docs.docker.com/compose/install/)

#### Steps
1. Start Docker daemon
2. cd to repo's root folder and run
    
    `go test -failfast -v -timeout 15m github.com/celer-network/sgn/test/e2e/multinode`
