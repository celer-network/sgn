## sgnops update-commission-rate decrease-rate

Decrease commission rate

### Synopsis

Decrease commission rate

```
sgnops update-commission-rate decrease-rate [flags]
```

### Options

```
      --add-lock-time string   (optional) additional rate lock period in unit of ETH block number
  -h, --help                   help for decrease-rate
      --rate string            Commission rate in unit of 0.01% (e.g., 120 is 1.2%)
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

* [sgnops update-commission-rate](sgnops_update-commission-rate.md)	 - Update commission rate

###### Auto generated by spf13/cobra
