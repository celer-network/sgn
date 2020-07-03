# SGN

[![CircleCI](https://circleci.com/gh/celer-network/sgn/tree/master.svg?style=svg)](https://circleci.com/gh/celer-network/sgn/tree/master)

## Running a validator on Ropsten

Pick a machine that has a minimum of 1gb memory (2gb recommended). Make sure you have go (version 1.13+), gcc, make, git and libleveldb-dev installed. Make sure you set \$GOBIN and add \$GOBIN to \$PATH. Eg:
`export GOBIN=~/go/bin;export GOPATH=~/go;export PATH=$PATH:\$GOBIN`

1. `git clone https://github.com/celer-network/sgn`
2. `cd sgn`
3. `WITH_CLEVELDB=yes make install`
4. `sgnd init [validator name] --chain-id sgnchain`
5. `sgncli config chain-id sgnchain; sgncli config output json; sgncli config indent true; sgncli config trust-node true`
6. `cd networks/ropsten`
7. `make copy-config`
8. `sgnd tendermint show-validator` to get the SGN validator public key
9. `sgncli keys add [key name] --keyring-backend=file` to add an account. You need at least one for the SGN operator / transactor. You can create more for multiple transactors. Note that the current implementation requires the same passphrase across all accounts. Use `sgncli keys list --keyring-backend=file` to get the list of SGN accounts info.
10. Prepare an Ethereum node endpoint. You can use services like Infura, Alchemy or your own node.
11. Prepare an Ethereum keystore and passphrase for the validator (eg. `geth account new`). Customize `config.json` with your Ethereum keystore path and passphrase, Ethereum full node instance url, SGN operator address, validator public key, list of transactors, and transactor/operator passphrase. Currently, operator and transactors share the same passphrase.
12. `sgnd start` to start the validator.
13. (optional) `sgncli gateway --laddr tcp://0.0.0.0:1317` to run a gateway.

Ropsten Mock Celer Token Address: 0xb37f671dFc6C7c03462C76313Ec1a35b0c0A76d5

In case that your local state is corrupted, you can try to reset state by running `sgnd unsafe-reset-all`.

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
- sgn nodeN log path: docker-volumes/nodeN/sgn/sgnd.log

### Manual Tests

### Setup

1. `cp ./test/config/local_config.json ./config.json`
1. start docker geth container `docker-compose up geth`
1. `WITH_CLEVELDB=yes make install-all`
1. `sgntest deploy`
1. `sgntest osp`
1. `sgnd start`

#### Test Guard

1. `curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seqNum": "10" }'`
1. `curl -X POST http://127.0.0.1:1317/intendSettle -d '{ "seqNum": "9" }'`

#### Test Upgrade

1. `sgncli tx govern submit-proposal software-upgrade test --title "upgrade test" --description "upgrade test" --deposit 10 --upgrade-height 100 --from jack --keyring-backend file`
1. `sgncli tx govern vote 1 yes --from jack --keyring-backend file`
1. Add upgrade handler to app.go, after the chain halts

   ```
   app.upgradeKeeper.SetUpgradeHandler("tesy", func(ctx sdk.Context, plan upgrade.Plan) {
      // upgrade changes here
      log.Infof("upgrade to tesy")
   })
   ```

1. Restart the chain `sgnd start`
