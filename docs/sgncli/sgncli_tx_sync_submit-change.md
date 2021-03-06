## sgncli tx sync submit-change

Submit a change

### Synopsis

Submit a change along with type and data.

Example:
$ <appcli> tx sync submit-change --type="sync_block" --data="My awesome change"

```
sgncli tx sync submit-change [flags]
```

### Options

```
      --blknum string   mainchain block number of change
      --data string     data of change
  -h, --help            help for submit-change
      --indent          Add indent to JSON response
      --trust-node      Trust connected full node (don't verify proofs for responses) (default true)
      --type string     type of change
```

### Options inherited from parent commands

```
      --config string     Path to SGN-specific configs (default "$HOME/.sgncli/config/sgn.toml")
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.sgncli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

### SEE ALSO

* [sgncli tx sync](sgncli_tx_sync.md)	 - Sync transactions subcommands

###### Auto generated by spf13/cobra
