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

1. `cp ./test/config/local_config.json ./config.json`
2. start docker geth container `docker-compose up geth`
3. `WITH_CLEVELDB=yes make install-all`
4. `sgntest deploy`
5. `sgntest osp`
6. `sgnd start`

#### Test Guard

1. `curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seqNum": "10" }'`
2. `curl -X POST http://127.0.0.1:1317/intendSettle -d '{ "seqNum": "9" }'`

#### Test Upgrade

1. `sgncli tx govern submit-proposal software-upgrade test --title "upgrade test" --description "upgrade test" --deposit 10 --upgrade-height 100`
2. `sgncli tx govern vote 1 yes`
3. Add upgrade handler to app.go, after the chain halts

```go
app.upgradeKeeper.SetUpgradeHandler("tesy", func(ctx sdk.Context, plan upgrade.Plan) {
// upgrade changes here
log.Infof("upgrade to tesy")
})
```

4. Restart the chain with `sgnd start`