## Sidechain Governance

### Prerequisite

Run `go run localnet.go -start -auto` to start testnet and auto config all nodes as validators.

### Example: update block reward

1. Query current block mining reward and submit change proposal:

Append args `--config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli` to the following commands.

```sh
sgncli query validator params
sgncli tx govern submit-proposal param-change data/param_change_proposal.json
sgncli query govern proposals
```

2. All nodes vote yes:

```sh
sgncli tx govern vote 1 yes --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli
sgncli tx govern vote 1 yes --config data/node1/sgncli/config/sgn.toml --home data/node1/sgncli
sgncli tx govern vote 1 yes --config data/node2/sgncli/config/sgn.toml --home data/node2/sgncli
```

3. Query proposal status and updated block mining reward after voting timeout (2 mins):

Append args `--config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli` to the following commands.

```sh
sgncli query govern proposal 1
sgncli query validator params
```