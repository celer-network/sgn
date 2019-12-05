package testing

import (
	"context"

	"github.com/celer-network/sgn/mainchain"
)

func DeployGuardContract(sgnParams *SGNParams) mainchain.Addr {
	ctx := context.Background()
	guardAddr, tx, _, err := mainchain.DeployGuard(EthClient.Auth, EthClient.Client, sgnParams.CelrAddr, sgnParams.BlameTimeout, sgnParams.MinValidatorNum, sgnParams.MinStakingPool, sgnParams.SidechainGoLiveTimeout)
	ChkErr(err, "failed to deploy Guard contract")
	WaitMinedWithChk(ctx, EthClient.Client, tx, 0, "Deploy Guard "+guardAddr.Hex())

	return guardAddr
}
