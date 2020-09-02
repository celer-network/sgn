package singlenode

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/celer-network/sgn/x/guard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGateway() []tc.Killable {
	p := &tc.SGNParams{
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1), // TODO: use a more practical value
		GovernVoteTimeout:      big.NewInt(1), // TODO: use a more practical value
		SlashTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MaxValidatorNum:        big.NewInt(11),
		MinStakingPool:         big.NewInt(0),
		AdvanceNoticePeriod:   big.NewInt(1), // TODO: use a more practical value
		SidechainGoLiveTimeout: big.NewInt(0),
		StartGateway:           true,
	}
	res := setupNewSGNEnv(p, "gateway")
	tc.SleepWithLog(10, "sgn being ready")

	return res
}

func TestE2EGateway(t *testing.T) {
	toKill := setupGateway()
	defer tc.TearDown(toKill)

	t.Run("e2e-gateway", func(t *testing.T) {
		t.Run("gatewayTest", gatewayTest)
	})
}

func gatewayTest(t *testing.T) {
	log.Info("=====================================================================")
	log.Info("======================== Test gateway ===========================")

	ctx := context.Background()
	cdc := app.MakeCodec()

	log.Info("Call subscribe on sgn contract...")
	amt, _ := new(big.Int).SetString("100000000000000000000", 10) // 100 CELR
	tx, err := tc.E2eProfile.CelrContract.Approve(tc.Client0.Auth, tc.E2eProfile.SGNAddr, amt)
	require.NoError(t, err, "failed to approve CELR on mainchain")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Approve CELR to SGN contract")

	tx, err = tc.SgnContract.Subscribe(tc.Client0.Auth, amt)
	require.NoError(t, err, "failed to subscribe on mainchain")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Subscribe on SGN contract")
	tc.SleepWithLog(10, "passing subscribe event block delay")

	resp, err := http.Get("http://127.0.0.1:1317/guard/subscription/" + tc.Client0.Address.Hex())
	require.NoError(t, err, "failed to query subscription from gateway")

	result, err := tc.ParseGatewayQueryResponse(resp, cdc)
	require.NoError(t, err, "failed to parse GatewayQueryResponse")

	var subscription guard.Subscription
	err = cdc.UnmarshalJSON(result, &subscription)
	require.NoError(t, err, "failed to unmarshal subscription JSON from gateway")

	log.Infoln("Query sgn about the subscription info:", subscription.String())
	expectedRes := fmt.Sprintf(`EthAddress: %s, Deposit: %d, Spend: %d`, mainchain.Addr2Hex(tc.Client0.Address), amt, 0) // defined in Subscription.String()
	assert.Equal(t, expectedRes, subscription.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
}
