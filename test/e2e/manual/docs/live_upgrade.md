## Live upgrade

The following steps perform a live upgrade on a local testnet. It does not require exporting or
migrating the state, thereby minimizing the downtime of the network. However, new validators will
need to use the old binary to replay the transactions during fast sync. Note that a live upgrade
will likely fail when the internals of Tendermint changes significantly. We will fall back to
exporting the state and migrating to a new network in that case.

### Prerequisite

Run `go run localnet.go -start -auto` to start testnet and auto config all nodes as validators.

### Propose and approve upgrade proposal

1. Query current block height and submit upgrade proposal:

Append args `--config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli` to the following commands.

```sh
sgncli query block
sgncli tx govern submit-proposal software-upgrade test --title "upgrade test" --description "upgrade test" --deposit 10 --upgrade-height [sidechain block height after more than 2 mins]
```

2. All nodes vote yes:

```sh
sgncli tx govern vote 1 yes --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli
sgncli tx govern vote 1 yes --config data/node1/sgncli/config/sgn.toml --home data/node1/sgncli
sgncli tx govern vote 1 yes --config data/node2/sgncli/config/sgn.toml --home data/node2/sgncli
```

3. Query proposal status after voting timeout (2 mins):

```sh
sgncli query govern proposal 1 --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli
```

4. Wait for the sidechain to halt at the proposed block height.

### Upgrade to the new version

1. Stop the containers:

```sh
go run localnet.go -stopall
```

2. Switch to the new code, add upgrade handler to `app.go`. Example: if upgrading from
[commit 1344abd](https://github.com/celer-network/sgn/tree/1344abd02183990f3f958fc3ae2b8ca148ee485f), use the following handler:

```go
app.upgradeKeeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan upgrade.Plan) {
    log.Info("upgrade to test")
    vparams := validator.Params{
        SyncerDuration: 10,
        MiningReward:   sdk.NewInt(10000000000000),
        PullerReward:   sdk.NewInt(500000000000),
        EpochLength:    3,
    }
    app.validatorKeeper.SetParams(ctx, vparams)
})
```

3. Rebuild the images:

```sh
go run localnet.go -rebuild
```

4. Start the new containers:

```sh
go run localnet.go -upall
```

### Note to new validator node after upgrade

New validator node who wants to join the sidechain after the upgrade should first run from genesis
using the old code to replay transactions before upgrade, then switch to the new code to replay
transactions after upgrade.

TODO: Check out state sync in Cosmos SDK v0.40.x (Stargate) to avoid replaying transactions from
block 0.

