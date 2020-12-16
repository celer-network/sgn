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
```

2. Run `go run localnet.go -start -auto` to start testnet and auto config all nodes as validators.

### Halt the nodes past a predetermined block height

1. Set the `halt-height` field in `../../../docker-volumes/nodeX/sgnd/config/app.toml` files to a
   block height greater than the intended height of the upgrade. For example, if the intended height
   is `300`, set to `301`.

2. Restart the nodes:

```sh
go run localnet.go -stopall
go run localnet.go -upall
```

3. Wait for the nodes to halt.

### Export the chain state

1. Stop the node processes:

```sh
go run localnet.go -stopall
```

2. Export the current state of the chain:

```sh
cd ../../../docker-volumes/node0
sgnd export --config ./sgncli/config/sgn.toml --home ./sgnd --for-zero-height --height <upgrade-height> > /tmp/sgntest_genesis_export.tmp
jq -S -M '' /tmp/sgntest_genesis_export.tmp > /tmp/sgntest_genesis_export.json # Sort and pretty-print the file
```

### Update the binary and migrate the genesis file

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
sgnd migrate [source_version] [target_version] [genesis_path] --chain-id sgntest-2 > sgntest-2_genesis.tmp
jq -S -M '' sgntest-2_genesis.tmp > sgntest-2_genesis.json # Sort and pretty-print the file
```

4. Rebuild the Docker images:

```sh
go run localnet.go -rebuild
```

5. Set `eth.monitor_start_block` to a past mainchain block number in the
   `data/nodeX/sgncli/config/sgn.toml` files.

### Reset the state and restart the nodes:

1. Replace the genesis files:

```sh
cp <path-to-new-genesis> ../../../docker-volumes/node0/sgnd/config/genesis.json # Repeat for all nodes
```

2. Set `eth.monitor_start_block` to a past mainchain block number in the
   `../../../docker-volumes/nodeX/sgncli/config/sgn.toml` files.

3. Set `halt-height` back to zero in `../../../docker-volumes/nodeX/sgnd/config/app.toml`.

4. Reset the local data on the nodes:

```sh
sgnd unsafe-reset-all --config ../../../docker-volumes/node0/sgncli/config/sgn.toml --home ../../../docker-volumes/node0/sgnd # Repeat for all nodes
```

5. Update `chain_id` in `../../../docker-volumes/nodeX/sgncli/config/sgn.toml` to the chain-id in the new genesis

6. Restart all nodes:

```sh
go run localnet.go -upall
```

7. Modify `../../../docker-volumes/nodeX/sgncli/config/sgn.toml` files and remove
   `eth.monitor_start_block`, so that future restarts will not try to monitor past events.
