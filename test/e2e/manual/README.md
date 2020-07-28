# Run Local Manual Tests

Follow instructions below to easily start a local testnet and play with multiple validator nodes on your local machine.

Start local testnet. `go run localnet.go -up`

`docker exec -ti sgnnode0 /bin/sh`

`sgnops init-candidate --commission-rate 1 --min-self-stake 1 --rate-lock-end-time 10000 --config config.json`

`sgncli query validator candidate 6a6d2a97da1c453a4e099e8054865a0a59728863 --home ./sgncli`

`sgnops delegate --candidate 6a6d2a97da1c453a4e099e8054865a0a59728863 --amount 10000 --config config.json`