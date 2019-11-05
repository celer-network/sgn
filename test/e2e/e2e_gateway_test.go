package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"

	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func setUpGateway() []tf.Killable {
	p := &SGNParams{
		blameTimeout:           big.NewInt(10),
		minValidatorNum:        big.NewInt(0),
		minStakingPool:         big.NewInt(0),
		sidechainGoLiveTimeout: big.NewInt(0),
		startGateway:           true,
	}
	res := setupNewSGNEnv(p, "gateway")
	sleepWithLog(20, "sgn being ready")

	return res
}

func TestE2EGateway(t *testing.T) {
	toKill := setUpGateway()
	defer tf.TearDown(toKill)

	t.Run("e2e-gateway", func(t *testing.T) {
		t.Run("gatewayTest", gatewayTest)
	})
}

func gatewayTest(t *testing.T) {
	// TODO: each test cases need a new and isolated sgn right now, which can't be run in parallel
	// t.Parallel()

	log.Info("=====================================================================")
	log.Info("======================== Test gateway ===========================")

	ctx := context.Background()
	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	ethAddress := tf.EthClient.Address
	guardContract := tf.EthClient.Guard
	transactor := tf.Transactor
	celrContract, err := mainchain.NewERC20(ctype.Hex2Addr(MockCelerAddr), conn)
	tf.ChkErr(err, "NewERC20 error")

	client1PrivKey, _ := crypto.HexToECDSA(client1Priv)
	client1Auth := bind.NewKeyedTransactor(client1PrivKey)
	client1Auth.GasPrice = big.NewInt(2e9) // 2Gwei

	// Call subscribe on guard contract
	log.Info("Call subscribe on guard contract...")
	amt := new(big.Int)
	amt.SetString("100000000000000000000", 10) // 100 CELR
	tx, err := celrContract.Approve(auth, ctype.Hex2Addr(GuardAddr), amt)
	tf.ChkErr(err, "failed to approve CELR to Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")
	tx, err = guardContract.Subscribe(auth, amt)
	tf.ChkErr(err, "failed to call subscribe of Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, maxBlockDiff+2, "Subscribe on Guard contract")

	msg := map[string]interface{}{
		"ethAddr": ethAddress.String(),
	}
	body, err := json.Marshal(msg)
	tf.ChkErr(err, "failed to marshal json")
	_, err = http.Post("https://127.0.0.1:1317/subscribe/subscribe", "application/json", bytes.NewBuffer(body))
	tf.ChkErr(err, "failed to post subscribe msg")

	resp, err := http.Get("https://127.0.0.1:1317/subscribe/subscription/" + ethAddress.String())
	tf.ChkErr(err, "failed to get subscription")
	body, err = ioutil.ReadAll(resp.Body)
	tf.ChkErr(err, "failed to read http response")
	var subscription subscribe.Subscription
	transactor.CliCtx.Codec.MustUnmarshalJSON(body, &subscription)
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	expectedRes := fmt.Sprintf(`Deposit: %d, Spend: %d`, amt, 0) // defined in Subscription.String()
	assert.Equal(t, subscription.String(), expectedRes, fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
}
