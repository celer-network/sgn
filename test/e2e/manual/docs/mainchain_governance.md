## Mainchain Governance

### Prerequisite

1. Run `go run localnet.go -start -auto` to start testnet and auto config all nodes as validators.

2. Checkout `sgnops gov --help` to see available commands. Use `sgnops gov get-param --help` to see the
parameter ID mapping.

### Example: update mainchain parameter

1. Query the current parameter `SlashTimeout`:

```sh
sgnops gov get-param 2 --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli
```

2. Create a proposal to change `SlashTimeout` to `30`:

```sh
sgnops gov create-param-proposal 2 30 --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli
```

3. Query latest submitted proposal:

```sh
sgnops gov get-param-proposal --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli # add --proposal-id for earlier proposals
```

4. All nodes vote yes:

```sh
sgnops gov vote-param-proposal 0 yes --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli
sgnops gov vote-param-proposal 0 yes --config data/node1/sgncli/config/sgn.toml --home data/node1/sgncli
sgnops gov vote-param-proposal 0 yes --config data/node2/sgncli/config/sgn.toml --home data/node2/sgncli
```

5. Query proposal and vote stats:

```sh
sgnops gov get-param-proposal --check-votes --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli
```

6. Confirm proposal after voting deadline:

```sh
sgnops gov confirm-param-proposal 0 --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli
```

7. Query updated `SlashTimeout` parameter:

```sh
sgnops gov get-param 2 --config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli
```
