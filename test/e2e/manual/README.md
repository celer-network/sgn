# Run Local Manual Tests

Follow instructions below to start a local testnet with three validator nodes on your local machine.

#### Table of contents

- [Add validators](#add-validators)
- [Sidechain governance](#sidechain-governance)
- [Mainchain governance](#mainchain-governance)

## Add validators

Run `go run localnet.go -start` to set up the docker test environment with three sgn nodes.

### Add node0 to become a validator
Append args `--config data/node0/config.json --home data/node0/sgncli` to following commands.

- `sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 300`
- `sgncli query validator candidate 00078b31fa8b29a76bce074b5ea0d515a6aeaee7`
- `sgnops delegate --candidate 00078b31fa8b29a76bce074b5ea0d515a6aeaee7 --amount 10000`
- `sgncli query validator validator sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk --trust-node`

### Add node1 to become a validator
Append args `--config data/node1/config.json --home data/node1/sgncli` to following commands.

#### Init node1 and self-delegate on mainchain
- `sgnops init-candidate --commission-rate 200 --min-self-stake 2000 --rate-lock-period 300`
- `sgncli query validator candidate 0015f5863ddc59ab6610d7b6d73b2eacd43e6b7e`
- `sgnops delegate --candidate 0015f5863ddc59ab6610d7b6d73b2eacd43e6b7e --amount 20000`
- `sgncli query validator validator sgn1egtta7su5jxjahtw56pe07qerz4lwvrlttac6y --trust-node`

### Add node2 to become a validator
Append args `--config data/node2/config.json --home data/node2/sgncli` to following commands.

- `sgnops init-candidate --commission-rate 120 --min-self-stake 3000 --rate-lock-period 300`
- `sgncli query validator candidate 00290a43e5b2b151d530845b2d5a818240bc7c70`
- `sgnops delegate --candidate 00290a43e5b2b151d530845b2d5a818240bc7c70 --amount 10000`
- `sgncli query validator validator sgn19q9usqmjcmx8vynynfl5tj5n2k22gc5f6wjvd7 --trust-node`

### Query all validators on sidechain
Append args `--config data/node0/config.json --home data/node0/sgncli` to following commands.
- `sgncli query validator validators`
- `sgncli query tendermint-validator-set --trust-node`

## Sidechain Governance

Run `go run localnet.go -start -auto` to start testnet and auto config all nodes as validators.

### Update block reward

#### Query current block mining reward and submit change proposal
Append args `--config data/node0/config.json --home data/node0/sgncli` to following commands.

- `sgncli query validator params`
- `sgncli tx govern submit-proposal param-change data/param_change_proposal.json`
- `sgncli query govern proposals`

#### All nodes vote yes
- `sgncli tx govern vote 1 yes --config data/node0/config.json --home data/node0/sgncli`
- `sgncli tx govern vote 1 yes --config data/node1/config.json --home data/node1/sgncli`
- `sgncli tx govern vote 1 yes --config data/node2/config.json --home data/node2/sgncli`

#### Query proposal status and updated block mining reward after voting timeout (2 mins)
Append args `--config data/node0/config.json --home data/node0/sgncli` to following commands.

- `sgncli query govern proposal 1`
- `sgncli query validator params`

### Upgrade sidechain

#### Query current block height and submit upgrade proposal
Append args `--config data/node0/config.json --home data/node0/sgncli` to following commands.

- `sgncli query block --trust-node`
- `sgncli tx govern submit-proposal software-upgrade test --title "upgrade test" --description "upgrade test" --deposit 10 --upgrade-height [sidechain block height after more than 2 mins]`

#### All nodes vote yes, same as [above](#all-nodes-vote-yes)

#### Query proposal status after voting timeout (2 mins)
- `sgncli query govern proposal 1 --config data/node0/config.json --home data/node0/sgncli`

The sidechain will stop at the proposed block height

#### Upgrade to the new version
1. Stop containers: `go run localnet.go -stopall`
2. Switch to the new code, add upgrade handler to `app.go`. Example: if upgrade from [commit 1344abd](https://github.com/celer-network/sgn/tree/1344abd02183990f3f958fc3ae2b8ca148ee485f), use the following handler
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
3. Rebuild images: `go run localnet.go -rebuild`
4. Start new containers: `go run localnet.go -upall`

#### Note to new validator node after upgrade
New validator node who wants to join the sidechain after the upgrade should first run from genesis using the old code to replay transactions before upgrade, then switch to the new code to replay transactions after upgrade.

## Mainchain Governance

Checkout `sgnops gov --help` to see available commands.

### Example: update mainchain parameter
1. **Query current param.** `sgnops gov get-param --help` to see the `param-id` mapping, then query current `SlashTimeout` value:

    `sgnops gov get-param 2 --config data/node0/config.json --home data/node0/sgncli`

2. **Create param change proposal.** Propose changing `SlashTimeout` to `30`:

   `sgnops gov create-param-proposal 2 30 --config data/node0/config.json --home data/node0/sgncli`

3. **Query submitted proposal** Get the latest proposal. If query previous proposal, add `--proposal-id` flag.

    `sgnops gov get-param-proposal --config data/node0/config.json --home data/node0/sgncli`

4. **All node vote yes**

    `sgnops gov vote-param-proposal 0 yes --config data/node0/config.json --home data/node0/sgncli`
 
    `sgnops gov vote-param-proposal 0 yes --config data/node1/config.json --home data/node1/sgncli`
 
    `sgnops gov vote-param-proposal 0 yes --config data/node2/config.json --home data/node2/sgncli`

5. **Query proposal and vote stats**

    `sgnops gov get-param-proposal --check-votes --config data/node0/config.json --home data/node0/sgncli`

6. **Confirm proposal after voting deadline**

    `sgnops gov confirm-param-proposal 0 --config data/node0/config.json --home data/node0/sgncli`

7. **Query updated param** Get current `SlashTimeout` value:

    `sgnops gov get-param 2 --config data/node0/config.json --home data/node0/sgncli`