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
	tc "github.com/celer-network/sgn/test/common"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setUpGateway() []tc.Killable {
	p := &tc.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MinStakingPool:         big.NewInt(0),
		SidechainGoLiveTimeout: big.NewInt(0),
		StartGateway:           true,
		CelrAddr:               tc.E2eProfile.CelrAddr,
		MaxValidatorNum:        big.NewInt(11),
	}
	res := setupNewSGNEnv(p, "gateway")
	tc.SleepWithLog(10, "sgn being ready")

	return res
}

func TestE2EGateway(t *testing.T) {
	toKill := setUpGateway()
	defer tc.TearDown(toKill)

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
	conn := tc.DefaultTestEthClient.Client
	auth := tc.DefaultTestEthClient.Auth
	ethAddress := tc.DefaultTestEthClient.Address
	guardContract := tc.DefaultTestEthClient.Guard
	transactor := tc.NewTransactor(
		t,
		CLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[1],
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
	)
	Client1PrivKey, err := tc.GetEthPrivateKey(tc.ValEthKs[1])
	tc.ChkTestErr(t, err, "failed to get client 1 private key")
	client1Auth := bind.NewKeyedTransactor(Client1PrivKey)
	client1Auth.GasPrice = big.NewInt(2e9) // 2Gwei

	log.Info("Call subscribe on guard contract...")
	amt, _ := new(big.Int).SetString("100000000000000000000", 10) // 100 CELR
	tx, err := tc.E2eProfile.CelrContract.Approve(auth, tc.E2eProfile.GuardAddr, amt)
	tc.ChkTestErr(t, err, "failed to approve CELR on mainchain")
	tc.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")

	tx, err = guardContract.Subscribe(auth, amt)
	tc.ChkTestErr(t, err, "failed to subscribe on mainchain")
	tc.WaitMinedWithChk(ctx, conn, tx, tc.BlockDelay, "Subscribe on Guard contract")
	tc.SleepWithLog(10, "passing subscribe event block delay")

	msg := map[string]interface{}{
		"ethAddr": ethAddress.Hex(),
	}
	body, _ := json.Marshal(msg)
	_, err = http.Post("http://127.0.0.1:1317/subscribe/subscribe", "application/json", bytes.NewBuffer(body))
	tc.ChkTestErr(t, err, "failed to post subscribe msg to gateway")
	tc.SleepWithLog(10, "sgn syncing Subscribe balance from mainchain")

	resp, err := http.Get("http://127.0.0.1:1317/subscribe/subscription/" + ethAddress.Hex())
	tc.ChkTestErr(t, err, "failed to query subscription from gateway")

	result, err := tc.ParseGatewayQueryResponse(resp, transactor.CliCtx.Codec)
	tc.ChkTestErr(t, err, "failed to parse GatewayQueryResponse")
	var subscription subscribe.Subscription
	err = transactor.CliCtx.Codec.UnmarshalJSON(result, &subscription)
	tc.ChkTestErr(t, err, "failed to unmarshal subscription JSON from gateway")
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	expectedRes := fmt.Sprintf(`Deposit: %d, Spend: %d`, amt, 0) // defined in Subscription.String()
	assert.Equal(t, expectedRes, subscription.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
}
