## Local Testing

### Multinode E2E Tests

#### Requirements

- Install [docker](https://docs.docker.com/install/)
- Install [docker-compose](https://docs.docker.com/compose/install/)

#### Steps

1. Start Docker daemon
2. `cd` into the root folder of the repo and run the following command (`sudo` may be required)

```sh
go test -failfast -v -timeout 30m github.com/celer-network/sgn/test/e2e/multinode
```

#### Test Logs

- geth log path: docker-volumes/geth-env/geth.log
- sgn nodeN log path: docker-volumes/nodeN/sgn/sgnd.log

### [Multinode Manual Tests](./e2e/manual/README.md)

### SingleNode Manual Tests

#### Setup

From the root folder of the repo:

```sh
WITH_CLEVELDB=yes make install-all
make reset-test-data
make localnet-start-geth
sgnops deploy --contract all
sgnd start 2>&1 | tee sgnd.log
```

#### Test Validator and Delegator

Note: wait for a few seconds between steps:

```sh
sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 300
sgncli tx validator set-transactors
sgncli query validator candidate 00078b31fa8b29a76bce074b5ea0d515a6aeaee7
sgnops delegate --candidate 00078b31fa8b29a76bce074b5ea0d515a6aeaee7 --amount 10000
sgncli query validator validator sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk
sgncli tx validator withdraw-reward 00078b31fa8b29a76bce074b5ea0d515a6aeaee7
```

#### Test Guard

```sh
sgncli gateway --laddr tcp://0.0.0.0:1318
sgnops guard-test --sgn-gateway http://127.0.0.1:1318
curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seq_num": "10" }' # should succeed
curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seq_num": "12" }' # should succeed
curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seq_num": "11" }' # should fail
curl -X POST http://127.0.0.1:1317/intendSettle -d '{ "seq_num": "9" }' # should success, look for guard tx in sgnd.log
curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seq_num": "15" }' # should fail
```

#### Test Live Upgrade

1. Submit and vote yes for the upgrade proposal:

```sh
sgncli tx govern submit-proposal software-upgrade test --title "upgrade test" --description "upgrade test" --deposit 10 --upgrade-height 100
sgncli tx govern vote 1 yes
```

2. Add upgrade handler to `app.go` after the chain halts:

```go
app.upgradeKeeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan upgrade.Plan) {
    // upgrade changes here
    log.Infof("upgrade to test")
})
```

3. Restart the chain with `sgnd start`