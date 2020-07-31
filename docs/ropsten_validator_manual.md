# Ropsten Validator Manual

1. Pick a Linux machine with a minimum of 1GB RAM (2GB recommended). Make sure you have go (version 1.14+), gcc, make, git and libleveldb-dev installed.

2. Make sure you set \$GOBIN and add \$GOBIN to \$PATH. Eg:

```shellscript
export GOBIN=~/go/bin;export GOPATH=~/go;export PATH=$PATH:\$GOBIN
```

Your actual paths might be different.

3. Clone the repository and install the `sgnd`, `sgncli` and `sgnops` binaries:

```shellscript
git clone https://github.com/celer-network/sgn
cd sgn
WITH_CLEVELDB=yes make install
make install-ops
```

4. Initialize the validator node:

```shellscript
sgnd init <validator-name> --chain-id sgnchain
```

`<validator-name>` can be any name of your choice.

5. Initialize `config.toml` containing general Cosmos SDK configs:

```shellscript
sgncli config chain-id sgnchain
sgncli config output json
sgncli config indent true
sgncli config trust-node true
```

6. Initialize `config.json` containing SGN specific configs:

```shellscript
cd networks/ropsten
make copy-config
```

7. Get the validator public key:

```
sgnd tendermint show-validator
```

Make a note of the output.

8. Add an SGN account key:

```
sgncli keys add <key-name> --keyring-backend=file
```

You need at least one for the SGN operator / transactor. You can create more for multiple transactors. Note that the current implementation requires the **same passphrase** for all accounts.

To get the list of accounts, run:

```
sgncli keys list --keyring-backend=file
```

9. Obtain a Ropsten Ethereum endpoint URL from [Infura](https://infura.io/). You may also use paid
   services like [Alchemy](https://alchemyapi.io/) or run your own node.

10. Prepare an Ethereum keystore for the validator. Eg.:

```shellscript
geth account new --lightkdf --keystore <path-to-keystore-folder>
```

11. Fill in the placeholders in `config.json`.
    | Field | Description |
    | ----- | ----------- |
    | ETHEREUM_GATEWAY_URL | The Ethereum gateway URL obtained from step 9 |
    | KEYSTORE_PATH | The path to the keystore file in step 10 |
    | KEYSTORE_PASSPHRASE | The passphrase to the keystore |
    | OPERATOR_ADDRESS | The cosmos-prefixed address obtained in step 8 |
    | TRANSACTOR_PASSPHRASE | The passphrase you typed in step 8 |
    | VALIDATOR_PUBKEY | The validator public key obtained in step 7 |
    | TRANSACTOR_ADDRESS | Reuse the operator address if you only created one account, or fill in multiple transactor accounts |

12. Start the validator:

```shellscript
sgnd start
```

It will take a while to sync the node.

13. Initialize the candidate status for your validator node:

```shellscript
sgnops init-candidate --commission-rate 1 --min-self-stake 1 --rate-lock-period 10000
```

It will take a while to complete the transactions on Ropsten.

14. Delegate 10000 Ropsten CELR to your candidate, which is the minimum amount required for it to
    become a validator:

```shellscript
sgnops delegate --candidate <candidate-eth-address> --amount 10000
```

`<candidate-eth-address>` is the ETH address obtained in step 10. It will take a while to complete
the transactions on Ropsten.

15. (Optional) In another terminal window, start an SGN gateway server:

```shellscript
sgncli gateway --laddr tcp://0.0.0.0:1317` to run a gateway.
```

16. Verify your candidate status at `http://54.218.106.24:8000/#/dpos`

17. (Optional) You can withdraw your self-stake and unbond your validator candidate by running:

```shellscript
sgnops intend-withdraw --candidate <candidate-eth-address> --amount 10000
```

After 2 hours, confirm the unbonded status and the withdrawal of your stake:

```shellscript
sgnops confirm-unbonded-candidate --candidate <candidate-eth-address>
sgnops confirm-withdraw --candidate <candidate-eth-address>
```

Each command will take a while to complete the transactions on Ropsten.

18. In case your local state is corrupted, you can try to reset the state by running:

```shellscript
sgnd unsafe-reset-all
```

**Please contact us before doing this**.
