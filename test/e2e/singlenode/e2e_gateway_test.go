package singlenode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setUpGateway() []tf.Killable {
	p := &tf.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MinStakingPool:         big.NewInt(0),
		SidechainGoLiveTimeout: big.NewInt(0),
		StartGateway:           true,
	}
	res := setupNewSGNEnv(p, "gateway")
	tf.SleepWithLog(10, "sgn being ready")

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
	conn := tf.DefaultTestEthClient.Client
	auth := tf.DefaultTestEthClient.Auth
	ethAddress := tf.DefaultTestEthClient.Address
	guardContract := tf.DefaultTestEthClient.Guard
	transactor := tf.NewTransactor(
		t,
		viper.GetString(common.FlagCLIHome),
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[0],
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
	)
	Client1PrivKey, _ := crypto.HexToECDSA(tf.Client1Priv)
	client1Auth := bind.NewKeyedTransactor(Client1PrivKey)
	client1Auth.GasPrice = big.NewInt(2e9) // 2Gwei

	log.Info("Call subscribe on guard contract...")
	amt, _ := new(big.Int).SetString("100000000000000000000", 10) // 100 CELR
	tx, err := tf.E2eProfile.CelrContract.Approve(auth, tf.E2eProfile.GuardAddr, amt)
	tf.ChkTestErr(t, err, "failed to approve CELR on mainchain")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")

	tx, err = guardContract.Subscribe(auth, amt)
	tf.ChkTestErr(t, err, "failed to subscribe on mainchain")
	tf.WaitMinedWithChk(ctx, conn, tx, tf.BlockDelay, "Subscribe on Guard contract")

	msg := map[string]interface{}{
		"ethAddr": ethAddress.Hex(),
	}
	body, _ := json.Marshal(msg)
	_, err = http.Post("http://127.0.0.1:1317/subscribe/subscribe", "application/json", bytes.NewBuffer(body))
	tf.ChkTestErr(t, err, "failed to post subscribe msg to gateway")
	tf.SleepWithLog(10, "sgn syncing Subscribe balance from mainchain")

	resp, err := http.Get("http://127.0.0.1:1317/subscribe/subscription/" + ethAddress.Hex())
	tf.ChkTestErr(t, err, "failed to query subscription from gateway")

	result, err := tf.ParseGatewayQueryResponse(resp, transactor.CliCtx.Codec)
	tf.ChkTestErr(t, err, "failed to parse GatewayQueryResponse")
	var subscription subscribe.Subscription
	err = transactor.CliCtx.Codec.UnmarshalJSON(result, &subscription)
	tf.ChkTestErr(t, err, "failed to unmarshal subscription JSON from gateway")
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	expectedRes := fmt.Sprintf(`Deposit: %d, Spend: %d`, amt, 0) // defined in Subscription.String()
	assert.Equal(t, expectedRes, subscription.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
}
