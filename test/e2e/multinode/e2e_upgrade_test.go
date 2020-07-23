package multinode

import (
	"math/big"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/testing/common"
	govtypes "github.com/celer-network/sgn/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setUpUpgrade() {
	log.Infoln("Set up new sgn env")
	setupNewSGNEnv(nil)
	tc.SleepWithLog(10, "sgn syncing")
}

func TestE2EUpgrade(t *testing.T) {
	setUpUpgrade()

	t.Run("e2e-upgrade", func(t *testing.T) {
		t.Run("upgradeTest", upgradeTest)
	})
}

func upgradeTest(t *testing.T) {
	log.Info("=====================================================================")
	log.Info("======================== Test upgrade ===========================")

	transactor0 := tc.NewTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnOperators[0],
		tc.SgnPassphrase,
	)

	amt1 := big.NewInt(3000000000000000000)
	amt2 := big.NewInt(2000000000000000000)
	amt3 := big.NewInt(2000000000000000000)
	amts := []*big.Int{amt1, amt2, amt3}
	tc.AddValidators(t, transactor0, tc.ValEthKs[:], tc.SgnOperators[:], amts)

	upgradeHeight := int64(100)
	plan := upgrade.Plan{Name: "test", Height: upgradeHeight}
	content := govtypes.NewUpgradeProposal("Upgrade test", "Upgrade test", plan)
	submitProposalmsg := govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor0.Key.GetAddress())
	transactor0.AddTxMsg(submitProposalmsg)

	proposalID := uint64(1)
	byteVoteOption, _ := govtypes.VoteOptionFromString("Yes")
	voteMsg := govtypes.NewMsgVote(transactor0.Key.GetAddress(), proposalID, byteVoteOption)
	transactor0.AddTxMsg(voteMsg)

	height, err := rpc.GetChainHeight(transactor0.CliCtx)
	for err == nil {
		time.Sleep(viper.GetDuration(common.FlagSgnTimeoutCommit) * time.Second)
		prevHeight := height
		height, err = rpc.GetChainHeight(transactor0.CliCtx)

		if prevHeight == height {
			break
		}
	}

	tc.ChkTestErr(t, err, "failed to query block height")
	assert.Equal(t, height, upgradeHeight-1, "The chain should stop at upgrade height")
}
