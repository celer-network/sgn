package testing

import (
	"context"
	"math/big"
	"strings"

	"github.com/celer-network/cChannel-eth-go/deploy"
	"github.com/celer-network/cChannel-eth-go/ledger"
	"github.com/celer-network/sgn/mainchain"
)

func DeployLedgerContract() mainchain.Addr {
	ctx := context.Background()
	channelAddrBundle := deploy.DeployAll(EthClient.Auth, EthClient.Client, ctx, 0)

	// Disable channel deposit limit
	LogBlkNum(EthClient.Client)
	ledgerContract, err := ledger.NewCelerLedger(channelAddrBundle.CelerLedgerAddr, EthClient.Client)
	ChkErr(err, "failed to NewCelerLedger")
	tx, err := ledgerContract.DisableBalanceLimits(EthClient.Auth)
	ChkErr(err, "failed disable channel deposit limits")
	WaitMinedWithChk(ctx, EthClient.Client, tx, 0, "Disable balance limit")

	return channelAddrBundle.CelerLedgerAddr
}

func DeployERC20Contract() (mainchain.Addr, *mainchain.ERC20) {
	ctx := context.Background()
	initAmt := new(big.Int)
	initAmt.SetString("5"+strings.Repeat("0", 44), 10)
	erc20Addr, tx, erc20, err := mainchain.DeployERC20(EthClient.Auth, EthClient.Client, initAmt, "Celer", 18, "CELR")
	ChkErr(err, "failed to deploy ERC20")
	WaitMinedWithChk(ctx, EthClient.Client, tx, 0, "Deploy ERC20 "+erc20Addr.Hex())

	return erc20Addr, erc20
}

func DeployGuardContract(sgnParams *SGNParams) mainchain.Addr {
	ctx := context.Background()
	guardAddr, tx, _, err := mainchain.DeployGuard(EthClient.Auth, EthClient.Client, sgnParams.CelrAddr, sgnParams.BlameTimeout, sgnParams.MinValidatorNum, sgnParams.MinStakingPool, sgnParams.SidechainGoLiveTimeout)
	ChkErr(err, "failed to deploy Guard contract")
	WaitMinedWithChk(ctx, EthClient.Client, tx, 0, "Deploy Guard "+guardAddr.Hex())

	return guardAddr
}
