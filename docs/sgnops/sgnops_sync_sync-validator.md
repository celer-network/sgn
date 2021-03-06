## sgnops sync sync-validator

Sync validator info from mainchain

### Synopsis

Example:
$ <appcli> tx submit-change sync-validator --candidate="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961"

```
sgnops sync sync-validator [flags]
```

### Options

```
      --candidate string   Candidate address
  -h, --help               help for sync-validator
      --indent             Add indent to JSON response
      --trust-node         Trust connected full node (don't verify proofs for responses) (default true)
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

* [sgnops sync](sgnops_sync.md)	 - Sync a change from mainchain to sidechain

###### Auto generated by spf13/cobra
