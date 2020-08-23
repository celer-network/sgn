# Run Local Manual Tests

Follow instructions below to easily start a local testnet and play with multiple validator nodes on your local machine.

## Start local testnet.

`go run localnet.go -up`

## Configure node0 to become a validator

`docker exec -ti sgnnode0 /bin/sh`

### Add validator
`sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 10000`

`sgnops delegate --candidate 6a6d2a97da1c453a4e099e8054865a0a59728863 --amount 10000`

### Query validator
`sgncli query validator candidate 6a6d2a97da1c453a4e099e8054865a0a59728863`

`sgncli query validator validators`

`sgncli query validator validator sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk --trust-node`

`sgncli query tendermint-validator-set --trust-node`

## Configure node1 to become a validator

`docker exec -ti sgnnode1 /bin/sh`

### Add validator
`sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 10000`

`sgnops delegate --candidate ba756d65a1a03f07d205749f35e2406e4a8522ad --amount 10000`

### Query validator
`sgncli query validator candidate ba756d65a1a03f07d205749f35e2406e4a8522ad`

`sgncli query validator validators`

`sgncli query validator validator sgn1egtta7su5jxjahtw56pe07qerz4lwvrlttac6y --trust-node`

`sgncli query tendermint-validator-set --trust-node`

## Configure node2 to become a validator

`docker exec -ti sgnnode2 /bin/sh`

### Add validator
`sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 10000`

`sgnops delegate --candidate f25d8b54fad6e976eb9175659ae01481665a2254 --amount 10000`

### Query validator
`sgncli query validator candidate f25d8b54fad6e976eb9175659ae01481665a2254`

`sgncli query validator validators`

`sgncli query validator validator sgn19q9usqmjcmx8vynynfl5tj5n2k22gc5f6wjvd7 --trust-node`

`sgncli query tendermint-validator-set --trust-node`
