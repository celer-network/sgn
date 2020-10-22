## Local Testing

### Multinode E2E Tests

#### Requirements

- Install [docker](https://docs.docker.com/install/)
- Install [docker-compose](https://docs.docker.com/compose/install/)

#### Steps

1. Start Docker daemon
2. cd to repo's root folder and run the following command (`sudo` may be required)

```shellscript
go test -failfast -v -timeout 30m github.com/celer-network/sgn/test/e2e/multinode
```

#### Test Logs

- geth log path: docker-volumes/geth-env/geth.log
- sgn nodeN log path: docker-volumes/nodeN/sgn/sgnd.log

### [Multinode Manual Tests](./e2e/manual/README.md)

### SingleNode Manual Tests

#### Setup

Do following operations at repo root folder
1. `cp test/data/.sgncli/config/sgn_template.toml test/data/.sgncli/config/sgn.toml`
2. `WITH_CLEVELDB=yes make install-all`
3. `sgnd unsafe-reset-all`
4. `make localnet-start-geth`
5. `sgnops deploy`
6. `sgnd start 2>&1 | tee sgnd.log`

#### Test Validator and Delegator

Wait for a few seconds between steps
1. `sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 300`
2. `sgncli tx validator set-transactors`
3. `sgncli query validator candidate 00078b31fa8b29a76bce074b5ea0d515a6aeaee7`
4. `sgnops delegate --candidate 00078b31fa8b29a76bce074b5ea0d515a6aeaee7 --amount 10000`
5. `sgncli query validator validator sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk`
6. `sgncli tx validator withdraw-reward 00078b31fa8b29a76bce074b5ea0d515a6aeaee7`

#### Test Guard

1. `sgnops channel`
2. `curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seqNum": "10" }'`
3. `curl -X POST http://127.0.0.1:1317/intendSettle -d '{ "seqNum": "9" }'`

#### Test Upgrade

1. `sgncli tx govern submit-proposal software-upgrade test --title "upgrade test" --description "upgrade test" --deposit 10 --upgrade-height 100`
2. `sgncli tx govern vote 1 yes`
3. Add upgrade handler to app.go, after the chain halts

```go
app.upgradeKeeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan upgrade.Plan) {
// upgrade changes here
log.Infof("upgrade to test")
})
```

4. Restart the chain with `sgnd start`