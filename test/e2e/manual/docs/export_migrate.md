## Export state and migrate to a new network

The following steps perform a migration to a new local testnet, which inherits the chain state but
starts from block height 0 again. This procedure is required for major upgrades that cannot be
handled by a live upgrade.

### Prerequisite

1. Make sure the `sgnd` binary on the host machine is up-to-date:

```sh
cd ../../../
make install
cd test/e2e/manual
``

2. Run `go run localnet.go -start -auto` to start testnet and auto config all nodes as validators.

### Halt the chain at a predetermined block height

1. Set the `halt-time` field in `../../../docker-volumes/nodeX/sgnd/config/app.toml` files to a future block height:

2. Restart the nodes:

```sh
go run localnet.go -stopall
go run localnet.go -upall
```

3. Wait for the nodes to halt.

### Export chain state

1. Stop the node processes:

```sh
go run localnet.go -stopall
```

2. Export the current state of the chain:

```sh
sgnd export --config ../../../docker-volumes/node0/sgncli/config/sgn.toml --home ../../../docker-volumes/node0/sgnd --for-zero-height --height <last-commit-height> > /tmp/sgntest_genesis_export.json
```

### Update the binary

1. Make a backwards-incompatible change and implement the migration command if needed. (TODO: add
more details)

2. Rebuild the `sgnd` binary on the container host machine:

```sh
cd ../../../
make install
cd test/e2e/manual
```

3. With the new `sgnd` binary, migrate the exported genesis file.

```sh
sgnd migrate <new-version> sgntest_genesis_export.json --chain-id sgntest-2 > sgntest-2_genesis.json
```

4. Replace the genesis files:

```sh
cp <path-to-new-genesis> data/node0/sgnd/config/genesis.json # Repeat for all nodes
```

5. Set `eth.monitor_start_block` to a past mainchain block number in the
`data/nodeX/sgncli/config/sgn.toml` files.

6. Rebuild the Docker images:

```sh
go run localnet.go -rebuild
```

### Reset the state and restart:

1. Reset the local data on the nodes:

```sh
sgnd unsafe-reset-all --config ../../../docker-volumes/node0/sgncli/config/sgn.toml --home ../../../docker-volumes/node0/sgnd # Repeat for all nodes
```

2. Restart all nodes:

```sh
go run localnet.go -upall
```

3. Modify `../../../docker-volumes/nodeX/sgncli/config/sgn.toml` files and remove
`eth.monitor_start_block`, so that future restarts will not try to monitor past events.
