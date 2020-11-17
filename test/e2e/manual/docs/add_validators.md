## Add validators

### Prerequisite

Run `go run localnet.go -start` to set up the docker test environment with three sgn nodes.

### Add node0 to become a validator

Append args `--config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli` to the following
commands.

```sh
sgnops init-candidate --commission-rate 150 --min-self-stake 1000 --rate-lock-period 300
sgncli tx validator set-transactors # This is only needed for the seed validator (node0).
sgncli query validator candidate 00078b31fa8b29a76bce074b5ea0d515a6aeaee7
sgnops delegate --candidate 00078b31fa8b29a76bce074b5ea0d515a6aeaee7 --amount 10000
sgncli query validator validator sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk
```

### Add node1 to become a validator

Append args `--config data/node1/sgncli/config/sgn.toml --home data/node1/sgncli` to the following
commands.

```sh
sgnops init-candidate --commission-rate 200 --min-self-stake 2000 --rate-lock-period 300
sgncli query validator candidate 0015f5863ddc59ab6610d7b6d73b2eacd43e6b7e
sgnops delegate --candidate 0015f5863ddc59ab6610d7b6d73b2eacd43e6b7e --amount 20000
sgncli query validator validator sgn1egtta7su5jxjahtw56pe07qerz4lwvrlttac6y
```

### Add node2 to become a validator

Append args `--config data/node2/sgncli/config/sgn.toml --home data/node2/sgncli` to the following
commands.

```sh
sgnops init-candidate --commission-rate 120 --min-self-stake 3000 --rate-lock-period 300
sgncli query validator candidate 00290a43e5b2b151d530845b2d5a818240bc7c70
sgnops delegate --candidate 00290a43e5b2b151d530845b2d5a818240bc7c70 --amount 10000
sgncli query validator validator sgn19q9usqmjcmx8vynynfl5tj5n2k22gc5f6wjvd7
```

### Query all validators on sidechain

Append args `--config data/node0/sgncli/config/sgn.toml --home data/node0/sgncli` to the following
commands.

```sh
sgncli query validator validators
sgncli query tendermint-validator-set
```
