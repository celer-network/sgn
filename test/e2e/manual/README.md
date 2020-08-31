# Run Local Manual Tests

Follow instructions below to easily start a local testnet and play with three validator nodes on your local machine.

#### Table of contents
- [Start local testnet](#start-local-testnet)
- [Add validators](#add-validators)

## Start local testnet.

Run `go run localnet.go -up`. This would set up the docker test environment, including one geth instance for mainchain and three sgn nodes for sidechain. This script also prepares local testnet accounts with testnet tokens.

## Add validators

### Add node0 to become a validator
Append args `--config data/node0/config.json --home data/node0/sgncli` in all commands.

Init and self-delegate on mainchain
- `sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 300`
- `sgnops delegate --candidate 6a6d2a97da1c453a4e099e8054865a0a59728863 --amount 10000`

Query validator on sidechain
- `sgncli query validator candidate 6a6d2a97da1c453a4e099e8054865a0a59728863`
- `sgncli query validator validator sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk --trust-node`

### Add node1 to become a validator
Append args `--config data/node1/config.json --home data/node1/sgncli` in all commands.

Init and self-delegate on mainchain
- `sgnops init-candidate --commission-rate 200 --min-self-stake 2000 --rate-lock-period 300`
- `sgnops delegate --candidate ba756d65a1a03f07d205749f35e2406e4a8522ad --amount 20000`

Query validator on sidechain
- `sgncli query validator candidate ba756d65a1a03f07d205749f35e2406e4a8522ad`
- `sgncli query validator validator sgn1egtta7su5jxjahtw56pe07qerz4lwvrlttac6y --trust-node`

### Add node2 to become a validator
Append args `--config data/node2/config.json --home data/node2/sgncli` in all commands.

nit and self-delegate on mainchain
- `sgnops init-candidate --commission-rate 120 --min-self-stake 3000 --rate-lock-period 300`
- `sgnops delegate --candidate f25d8b54fad6e976eb9175659ae01481665a2254 --amount 10000`

Query validator on sidechain
- `sgncli query validator candidate f25d8b54fad6e976eb9175659ae01481665a2254`
- `sgncli query validator validator sgn19q9usqmjcmx8vynynfl5tj5n2k22gc5f6wjvd7 --trust-node`

### Query all validators on sidechain
Append args `--config data/node0/config.json --home data/node0/sgncli` in all commands.
- `sgncli query validator validators`
- `sgncli query tendermint-validator-set --trust-node`