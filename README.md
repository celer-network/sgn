# SGN

[![CircleCI](https://circleci.com/gh/celer-network/sgn/tree/master.svg?style=svg)](https://circleci.com/gh/celer-network/sgn/tree/master)

## Ropsten Test User Manual

1. Clone the repository and install the `sgnops` binary:

```shellscript
git clone https://github.com/celer-network/sgn
cd sgn
make install-ops
cd networks/ropsten
```

2. Obtain a Ropsten Ethereum endpoint URL from [Infura](https://infura.io/).
3. Fill in the `ETHEREUM_GATEWAY_URL` in `config.json`. You can leave the other placeholders unfilled.
4. Create two keystores with **empty passphrase** for testing purpose. Eg.:

```shellscript
geth account new --lightkdf --keystore <path-to-keystore-folder>
```

5. Join our [Discord](https://discord.gg/uGx4fjQ)
   server and ping us to obtain some Ropsten mock CELR tokens. You should also obtain Ropsten ETH from places like the MetaMask [faucet](https://faucet.metamask.io).
6. Send Ropsten ETH and CELR to `peer1` and `peer2`. Make sure `peer1` has at least 1 Ropsten CELR.
   You can do so by importing the keystore JSON files into MetaMask.
7. Start the local test server:

```shellscript
sgnops channel --peer1 <path-to-peer1-keystore> --peer2 <path-to-peer2-keystore> --gateway http://54.218.106.24:1317
```

The test program will open a Celer Channel between the peers and subscribe `peer1` to the SGN. Once the bootstrap process is done, you will see “Starting RPC HTTP server on 127.0.0.1:1317”

8. Check if the subscription succeeded:

```shellscript
curl http://54.218.106.24:1317/guard/subscription/<peer1-address>
```

If not, you can run:

```shellscript
curl -X POST http://54.218.106.24:1317/guard/subscribe -d '{ "ethAddr": "<peer1-address>", "amount": "1000000000000000000" }'
```

to retry manually.

9. In the following command, the two peers co-sign a new state with sequence number 10. `peer1` then sends the state to SGN to be guarded.

```shellscript
curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seqNum": "10" }'
```

10. Check if the subscription succeeded:

```shellscript
curl http://54.218.106.24:1317/guard/request/<channel-id>/<peer1-address>
```

11. Now let `peer2` try to maliciously settle the channel with sequence number 9:

```shellscript
curl -X POST http://127.0.0.1:1317/intendSettle -d '{ "seqNum": "9" }'
```

12. Check if the SGN guards the channel successfully:

```shellscript
curl http://127.0.0.1:1317/channelInfo
```

If so, `seqNum` should be 10. Note that it can take a few minutes for this to happen.

## Ropsten Validator Manual

1. Pick a Linux machine with a minimum of 1GB RAM (2GB recommended).Make sure you have go (version 1.13+), gcc, make, git and libleveldb-dev installed.

2. Make sure you set \$GOBIN and add \$GOBIN to \$PATH. Eg:

```shellscript
export GOBIN=~/go/bin;export GOPATH=~/go;export PATH=$PATH:\$GOBIN
```

Your actual paths might be different.

3. Clone the repository and install the `sgnd` and `sgncli` binaries:

```shellscript
git clone https://github.com/celer-network/sgn
cd sgn
WITH_CLEVELDB=yes make install
```

4. Initialize the validator node:

```shellscript
sgnd init <validator-name> --chain-id sgnchain
```

`<validator-name>` can be any name of your choice.

5.

```shellscript
sgncli config chain-id sgnchain; sgncli config output json; sgncli config indent true; sgncli config trust-node true
```

6. Initialize the validator config file:

```shellscript
cd networks/ropsten
make copy-config
```

7. Get the validator public key:

```
sgnd tendermint show-validator
```

Make a note of the output.

8. Add an SGN account key:

```
sgncli keys add <key-name> --keyring-backend=file
```

You need at least one for the SGN operator / transactor. You can create more for multiple transactors. Note that the current implementation requires the **same passphrase** for all accounts.

To get the list of accounts, run:

```
sgncli keys list --keyring-backend=file
```

9. Obtain a Ropsten Ethereum endpoint URL from [Infura](https://infura.io/). You may also use paid
   services like [Alchemy](https://alchemyapi.io/) or run your own node.

10. Prepare an Ethereum keystore for the validator. Eg.:

```shellscript
geth account new --lightkdf --keystore <path-to-keystore-folder>
```

11. Fill in the placeholders in `config.json`.
    | Field | Description |
    | ----- | ----------- |
    | ETHEREUM_GATEWAY_URL | The Ethereum gateway URL obtained from step 9 |
    | KEYSTORE_PATH | The path to the keystore file in step 10 |
    | KEYSTORE_PASSPHRASE | The passphrase to the keystore |
    | OPERATOR_ADDRESS | The cosmos-prefixed address obtained in step 8 |
    | TRANSACTOR_PASSPHRASE | The passphrase you typed in step 8 |
    | VALIDATOR_PUBKEY | The validator public key obtained in step 7 |
    | TRANSACTOR_ADDRESS | Reuse the operator address if you only created one account, or fill in multiple transactor accounts |

12. Start the validator:

```shellscript
sgnd start
```

13. (Optional) In another terminal window, start an SGN gateway server:

```shellscript
sgncli gateway --laddr tcp://0.0.0.0:1317` to run a gateway.
```

14. In case your local state is corrupted, you can try to reset the state by running:

```shellscript
sgnd unsafe-reset-all
```

**Please contact us before doing this**.

## Local Testing

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

```

```
