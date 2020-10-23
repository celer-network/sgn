package singlenode

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupValidator() []tc.Killable {
	p := &tc.SGNParams{
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1),
		GovernVoteTimeout:      big.NewInt(1),
		SlashTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(1),
		MaxValidatorNum:        big.NewInt(11),
		MinStakingPool:         big.NewInt(1),
		AdvanceNoticePeriod:    big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
	}
	res := setupNewSGNEnv(p, "validator")
	tc.SleepWithLog(10, "sgn being ready")

	return res
}

func TestE2EValidator(t *testing.T) {
	toKill := setupValidator()
	defer tc.TearDown(toKill)

	t.Run("e2e-validator", func(t *testing.T) {
		t.Run("validatorTest", validatorTest)
	})
}

func validatorTest(t *testing.T) {
	log.Info("===================================================================")
	log.Info("======================== Test validator ===========================")

	transactor := tc.NewTestTransactor(
		t,
		CLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[0],
		viper.GetString(common.FlagSgnPassphrase),
	)
	vAmt := big.NewInt(1000000000000000000) // 1 CELR
	dAmts := []*big.Int{
		big.NewInt(2000000000000000000), // 2 CELR
		big.NewInt(2000000000000000000), // 2 CELR
		big.NewInt(4000000000000000000), // 4 CELR
		big.NewInt(1000000000000000000), // 1 CELR
	}
	miningPool := new(big.Int)
	miningPool.SetString("1"+strings.Repeat("0", 20), 10)

	vEthAddr, vAuth, err := tc.GetAuth(tc.ValEthKs[0])
	log.Infof("validator eth address %x", vEthAddr)
	require.NoError(t, err, "failed to get validator auth")

	log.Infof("add mining reward pool in contract %x", tc.E2eProfile.DPoSAddr)
	tx, err := tc.E2eProfile.CelrContract.Approve(vAuth, tc.E2eProfile.DPoSAddr, miningPool)
	require.NoError(t, err, "failed to approve CELR to DPoS contract")
	tc.WaitMinedWithChk(
		context.Background(), tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Approve CELR to DPoS contract")
	_, err = tc.DposContract.ContributeToMiningPool(vAuth, miningPool)
	require.NoError(t, err, "failed to call ContributeToMiningPool of DPoS contract")

	tc.AddCandidateWithStake(
		t, transactor, vEthAddr, vAuth, tc.ValAccounts[0], vAmt,
		big.NewInt(100000), /* minSelfStake */
		big.NewInt(200) /* commission rate 2% */, big.NewInt(10000) /* rate lock end time */, true)
	tc.CheckValidatorNum(t, transactor, 1)

	log.Info("add delegators ...")
	for i := 0; i < len(tc.DelEthKs); i++ {
		_, dAuth, err2 := tc.GetAuth(tc.DelEthKs[i])
		require.NoError(t, err2, "failed to get delegator auth")
		go tc.DelegateStake(dAuth, vEthAddr, dAmts[i])
	}
	for i := 0; i < len(tc.DelEthKs); i++ {
		tc.CheckDelegator(t, transactor, vEthAddr, mainchain.Hex2Addr(tc.DelEthAddrs[i]), dAmts[i])
	}

	time.Sleep(10 * time.Second)
	log.Info("check rewards ...")
	for i := 0; i < len(tc.DelEthKs); i++ {
		reward, err2 := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, tc.DelEthAddrs[i])
		require.NoError(t, err2, "failed to query reward on sgn")
		log.Infoln("sgn reward", reward.String())
		assert.True(t, reward.MiningReward.IsPositive(), "Minging reward should be larger than 0")
	}

	log.Info("init withdraw rewards ...")
	for i := 0; i < len(tc.DelEthKs); i++ {
		msgWithdrawReward := validator.NewMsgWithdrawReward(tc.DelEthAddrs[i], transactor.Key.GetAddress())
		transactor.AddTxMsg(msgWithdrawReward)
	}

	time.Sleep(5 * time.Second)
	log.Info("withdraw rewards on mainchain ...")
	var lastTx *types.Transaction
	for i := 0; i < len(tc.DelEthKs); i++ {
		reward, err2 := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, tc.DelEthAddrs[i])
		require.NoError(t, err2, "failed to query reward on sgn")
		_, dAuth, _ := tc.GetAuth(tc.DelEthKs[i])
		tx, err2 := tc.SgnContract.RedeemReward(dAuth, reward.GetRewardRequest())
		require.NoError(t, err2, "failed to redeem reward on mainchain")
		log.Infoln("redeem reward on mainchain, tx", tx.Hash().Hex())
		lastTx = tx
	}
	tc.WaitMinedWithChk(
		context.Background(), tc.EthClient, lastTx, tc.BlockDelay, tc.PollingInterval, "redeem reward on mainchain")

	for i := 0; i < len(tc.DelEthKs); i++ {
		r, err2 := tc.DposContract.RedeemedMiningReward(&bind.CallOpts{}, mainchain.Hex2Addr(tc.DelEthAddrs[i]))
		require.NoError(t, err2, "failed to query reward on mainchain")
		log.Infof("redeemed reward of delegator %s is %s", tc.DelEthAddrs[i], r)
		assert.True(t, r.Cmp(big.NewInt(0)) > 0, "Redeemed reward should be larger than 0")
	}
}
