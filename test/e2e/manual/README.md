# Run Local Manual Tests

Follow instructions below to easily start a local testnet and play with three validator nodes on your local machine.

## Start local testnet.

`go run localnet.go -up`

## Add validators

### Add node0 to become a validator
`docker exec -ti sgnnode0 /bin/sh`

#### Add validator
- `sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 300`
- `sgnops delegate --candidate 6a6d2a97da1c453a4e099e8054865a0a59728863 --amount 10000`

#### Confirm validator
- `sgncli query validator candidate 6a6d2a97da1c453a4e099e8054865a0a59728863`
- `sgncli query validator validator sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk --trust-node`

### Add node1 to become a validator
`docker exec -ti sgnnode1 /bin/sh`

#### Add validator
- `sgnops init-candidate --commission-rate 200 --min-self-stake 2000 --rate-lock-period 300`
- `sgnops delegate --candidate ba756d65a1a03f07d205749f35e2406e4a8522ad --amount 20000`

#### Confirm validator
- `sgncli query validator candidate ba756d65a1a03f07d205749f35e2406e4a8522ad`
- `sgncli query validator validator sgn1egtta7su5jxjahtw56pe07qerz4lwvrlttac6y --trust-node`

### add node2 to become a validator
`docker exec -ti sgnnode2 /bin/sh`

#### Add validator
- `sgnops init-candidate --commission-rate 120 --min-self-stake 3000 --rate-lock-period 300`
- `sgnops delegate --candidate f25d8b54fad6e976eb9175659ae01481665a2254 --amount 10000`

#### Confirm validator
- `sgncli query validator candidate f25d8b54fad6e976eb9175659ae01481665a2254`
- `sgncli query validator validator sgn19q9usqmjcmx8vynynfl5tj5n2k22gc5f6wjvd7 --trust-node`

### Query all validators
- `sgncli query validator validators`
- `sgncli query tendermint-validator-set --trust-node`