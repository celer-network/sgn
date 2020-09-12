# Run Local Manual Tests

Follow instructions below to start a local testnet with three validator nodes on your local machine.

#### Table of contents

- [Add validators](#add-validators)
- [Update param through governance](#update-param-through-governance)

## Add validators

Run `go run localnet.go -start` to set up the docker test environment with three sgn nodes.

### Add node0 to become a validator
Append args `--config data/node0/config.json --home data/node0/sgncli` to following commands.

- `sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 300`
- `sgncli query validator candidate 6a6d2a97da1c453a4e099e8054865a0a59728863`
- `sgnops delegate --candidate 6a6d2a97da1c453a4e099e8054865a0a59728863 --amount 10000`
- `sgncli query validator validator sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk --trust-node`

### Add node1 to become a validator
Append args `--config data/node1/config.json --home data/node1/sgncli` to following commands.

#### Init node1 and self-delegate on mainchain
- `sgnops init-candidate --commission-rate 200 --min-self-stake 2000 --rate-lock-period 300`
- `sgncli query validator candidate ba756d65a1a03f07d205749f35e2406e4a8522ad`
- `sgnops delegate --candidate ba756d65a1a03f07d205749f35e2406e4a8522ad --amount 20000`
- `sgncli query validator validator sgn1egtta7su5jxjahtw56pe07qerz4lwvrlttac6y --trust-node`

### Add node2 to become a validator
Append args `--config data/node2/config.json --home data/node2/sgncli` to following commands.

- `sgnops init-candidate --commission-rate 120 --min-self-stake 3000 --rate-lock-period 300`
- `sgncli query validator candidate f25d8b54fad6e976eb9175659ae01481665a2254`
- `sgnops delegate --candidate f25d8b54fad6e976eb9175659ae01481665a2254 --amount 10000`
- `sgncli query validator validator sgn19q9usqmjcmx8vynynfl5tj5n2k22gc5f6wjvd7 --trust-node`

### Query all validators on sidechain
Append args `--config data/node0/config.json --home data/node0/sgncli` to following commands.
- `sgncli query validator validators`
- `sgncli query tendermint-validator-set --trust-node`

## Update param through governance

Update block block reward through governance. 

Run `go run localnet.go -start -auto` to start testnet and auto config all nodes as validators.

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