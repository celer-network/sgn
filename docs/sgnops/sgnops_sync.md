## sgnops sync

Sync a change from mainchain to sidechain

### Synopsis

Sync a change from mainchain to sidechain

```
sgnops sync [flags]
```

### Options

```
  -h, --help   help for sync
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

* [sgnops](sgnops.md)	 - sgn ops utility
* [sgnops sync sync-delegator](sgnops_sync_sync-delegator.md)	 - Sync delegator info from mainchain
* [sgnops sync sync-subscription-balance](sgnops_sync_sync-subscription-balance.md)	 - Sync subscription balance info from mainchain
* [sgnops sync sync-update-sidechain-addr](sgnops_sync_sync-update-sidechain-addr.md)	 - Sync sidechain address from mainchain
* [sgnops sync sync-validator](sgnops_sync_sync-validator.md)	 - Sync validator info from mainchain

###### Auto generated by spf13/cobra
