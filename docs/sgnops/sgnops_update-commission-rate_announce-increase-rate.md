## sgnops update-commission-rate announce-increase-rate

Announce increase commission rate

### Synopsis

Announce increase commission rate

```
sgnops update-commission-rate announce-increase-rate [flags]
```

### Options

```
      --add-lock-time string   (optional) additional rate lock period
  -h, --help                   help for announce-increase-rate
      --rate string            Commission rate in unit of 0.01% (e.g., 120 is 1.2%)
```

### Options inherited from parent commands

```
      --config string     config path (default "./config.json")
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.sgncli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

### SEE ALSO

* [sgnops update-commission-rate](sgnops_update-commission-rate.md)	 - Update commission rate

###### Auto generated by spf13/cobra on 18-Aug-2020