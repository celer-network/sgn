# Ropsten Test User Manual

1. Clone the repository and install the `sgnops` binary:

```shellscript
git clone https://github.com/celer-network/sgn
cd sgn
make install-ops
cd networks/ropsten
```

2. Obtain a Ropsten Ethereum endpoint URL from [Infura](https://infura.io/).
3. Fill in the `ETHEREUM_GATEWAY_URL` in `config.json`. You can leave the other placeholders unfilled.
4. Create two keystores with **empty passphrase** for testing purpose. Eg.:

```shellscript
geth account new --lightkdf --keystore <path-to-keystore-folder>
```

5. Join our [Discord](https://discord.gg/uGx4fjQ)
   server and ping us to obtain some Ropsten mock CELR tokens. You should also obtain Ropsten ETH from places like the MetaMask [faucet](https://faucet.metamask.io).
6. Send Ropsten ETH and CELR to `peer1` and `peer2`. Make sure `peer1` has at least 1 Ropsten CELR.
   You can do so by importing the keystore JSON files into MetaMask.
7. Start the local test server:

```shellscript
sgnops channel --peer1 <path-to-peer1-keystore> --peer2 <path-to-peer2-keystore> --gateway http://54.218.106.24:1317
```

The test program will open a Celer Channel between the peers and subscribe `peer1` to the SGN. Once the bootstrap process is done, you will see “Starting RPC HTTP server on 127.0.0.1:1317”

8. Check if the subscription succeeded:

```shellscript
curl http://54.218.106.24:1317/guard/subscription/<peer1-address>
```

If not, you can run:

```shellscript
curl -X POST http://54.218.106.24:1317/guard/subscribe -d '{ "ethAddr": "<peer1-address>", "amount": "1000000000000000000" }'
```

to retry manually.

9. In the following command, the two peers co-sign a new state with sequence number 10. `peer1` then sends the state to SGN to be guarded.

```shellscript
curl -X POST http://127.0.0.1:1317/requestGuard -d '{ "seqNum": "10" }'
```

10. Check if the subscription succeeded:

```shellscript
curl http://54.218.106.24:1317/guard/request/<channel-id>/<peer1-address>
```

11. Now let `peer2` try to maliciously settle the channel with sequence number 9:

```shellscript
curl -X POST http://127.0.0.1:1317/intendSettle -d '{ "seqNum": "9" }'
```

12. Check if the SGN guards the channel successfully:

```shellscript
curl http://127.0.0.1:1317/channelInfo
```

If so, `seqNum` should be 10. Note that it can take a few minutes for this to happen.
