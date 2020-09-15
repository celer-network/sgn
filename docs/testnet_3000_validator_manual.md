# Testnet 3000 Validator Manual

**Note: This manual assumes familiarity with Unix command line and blockchain validator nodes.
Prior experience with Cosmos SDK will be helpful.**

1. Pick a Linux machine with a minimum of 1GB RAM (2GB recommended). Make sure you have go
   (version 1.14), gcc, make, git and libleveldb-dev installed. **NOTE: go1.15 seems to have an
   issue with the keyring library Cosmos SDK depends on. We are still investigating it.**

2. Make sure you set \$GOBIN and add \$GOBIN to \$PATH. Eg:

```shellscript
export GOBIN=~/go/bin;export GOPATH=~/go;export PATH=$PATH:\$GOBIN
```

Your actual paths might be different.

3. Clone the `sgn` repository and install the `sgnd`, `sgncli` and `sgnops` binaries:

```shellscript
git clone https://github.com/celer-network/sgn
cd sgn
git checkout master
WITH_CLEVELDB=yes make install-all
cd ..
```

4. Initialize the validator node:

```shellscript
sgnd init <validator-name> --chain-id sgn-testnet-3000
```

`<validator-name>` can be any name of your choice.

5. Clone the `sgn-testnets` repository:

```shellscript
git clone https://github.com/celer-network/sgn-testnets
cd sgn-testnets
git checkout master
cd 3000
```

You will notice a `config.json` file in the directory. Unless otherwise specified, all the `sgnd`
and `sgncli` commands in the following steps are run from the `3000` directory. If you would like to
run `sgncli` from a different directory, you can use the `--config` flag. Eg.
`sgncli --config <path-to-config.json> <cmd>`.

6. Copy the `genesis.json` and `config.toml` files to the validator node directory:

```shellscript
cp genesis.json config.toml $HOME/.sgnd/config
```

7. Fill out the `moniker` field in `$HOME/.sgnd/config/config.toml` with the name of your validator.

8. Get the validator public key:

```shellscript
sgnd tendermint show-validator
```

Make a note of the output.

9. Add an SGN account key:

```shellscript
sgncli keys add <key-name> --keyring-backend=file
```

You need at least one for the SGN validator / main transactor. You can create more for multiple
transactors. Note that the current implementation requires the **same passphrase** for all accounts.

To get the list of accounts, run:

```shellscript
sgncli keys list --keyring-backend=file
```

10. Obtain a Ropsten Ethereum endpoint URL from [Infura](https://infura.io/). You may also use paid
   services like [Alchemy](https://alchemyapi.io/) or run your own node.

11. Prepare an Ethereum keystore for the validator. Eg.:

```shellscript
geth account new --lightkdf --keystore <path-to-keystore-folder>
```

12. Join our [Discord](https://discord.gg/uGx4fjQ) server and ping us to obtain some Ropsten mock
    CELR tokens. You should also obtain a few Ropsten ETH from places like the MetaMask
    [faucet](https://faucet.metamask.io).

13. Send Ropsten ETH and CELR to the created ETH address. Make sure it has at least 10000 Ropsten
    CELR and a few Ropsten ETH for gas. You can import the keystore JSON file into MetaMask for ease
    of use.

14. Fill in the placeholders in `config.json`.
    | Field | Description |
    | ----- | ----------- |
    | ETHEREUM_GATEWAY_URL | The Ethereum gateway URL obtained from step 10 |
    | KEYSTORE_PATH | The path to the keystore file in step 11 |
    | KEYSTORE_PASSPHRASE | The passphrase to the keystore |
    | VALIDATOR_ACCOUNT_ADDRESS | The sgn-prefixed address obtained in step 9 |
    | TRANSACTOR_PASSPHRASE | The passphrase you typed in step 9 |
    | VALIDATOR_PUBKEY | The validator public key obtained in step 8 |
    | TRANSACTOR_ADDRESS | Reuse the validator account address if you only created one account, or fill in multiple transactor accounts |

15. Start the validator and redirect the output to a log file:

```shellscript
sgnd start > sgnd.log 2>&1
```

It will take a while to sync the node.

16. Currently, we maintain a whitelist of validators on the mainchain contract. Please report your
    validator ETH address in [Discord](https://discord.gg/uGx4fjQ) to get whitelisted. Note that the
    number of active validators is limited, so a slot is not guaranteed.

17. Initialize the candidate status for your validator node:

```shellscript
sgnops init-candidate --commission-rate 500 --min-self-stake 1000 --rate-lock-period 10000
```

It will take a while to complete the transactions on Ropsten.

18. Delegate 10000 Ropsten CELR to your candidate, which is the minimum amount required for it to
    become a validator:

```shellscript
sgnops delegate --candidate <candidate-eth-address> --amount 10000
```

`<candidate-eth-address>` is the ETH address obtained in step 11. It will take a while to complete
the transactions on Ropsten.

19. Note that it will take some time for the existing SGN validators to sync your new validator from
    the mainchain. After a while, verify your validator status:

```shellscript
sgncli query validator candidate <candidate-eth-address>
```

You should be able to see that your candidate has a `delegatedStake` of `10000000000000000000000`,
which is 10000 Ropsten CELR denominated in wei.

20. Verify your validator is in the SGN validator set:

```shellscript
sgncli query validator validator <validator-account-address> --trust-node
```

`<validator-account-address>` is the sgn-prefixed address obtained in step 9. You should see your
validator and its `identity` should equal to your `<candidate-eth-address>`. Make a note of the
`consensus_pubkey` - the address prefixed with `sgnvalconspub`.

21. Due to inaccurate Ethereum gas estimation, your validator might fail to claim its validator
status on mainchain. If your validator doesn't appear in the query, try claiming the status
manually:

```shellscript
sgnops claim-validator
```

If the command succeeds, wait for a while and retry they query again. If the command fails or your
validator still doesn't show up, please **contact us** on [Discord](https://discord.gg/uGx4fjQ).

22. Verify your validator is in the Tendermint validator set:

```shellscript
sgncli query tendermint-validator-set --trust-node
```

You should see an entry with `pub_key` matching the `consensus_pubkey` obtained in step 20.

You should also be able to see your validator on the dashboard at
`http://sgntest.celer.network/#/dpos`.

23. (Optional) You can run an SGN gateway server to serve HTTP requests. In another terminal window:

```shellscript
sgncli gateway --laddr tcp://0.0.0.0:1317` to run a gateway.
```

24. (Optional) You can withdraw your self-stake and unbond your validator candidate by running:

```shellscript
sgnops withdraw intend --candidate <candidate-eth-address> --amount 10000
```

After at least half an hour, confirm the unbonded status and the withdrawal of your stake:

```shellscript
sgnops confirm-unbonded-candidate --candidate <candidate-eth-address>
sgnops withdraw confirm --candidate <candidate-eth-address>
```

Each command will take a while to complete the transactions on Ropsten.

25. In case the local state of your validator is corrupted, you can try to reset the state by running:

```shellscript
sgnd unsafe-reset-all
```

Please **contact us** on [Discord](https://discord.gg/uGx4fjQ) before doing this.
