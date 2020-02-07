package testcommon

import (
	"context"
	"math/big"
	"strings"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/test/channel-eth-go/deploy"
	"github.com/celer-network/sgn/test/channel-eth-go/ledger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeployLedgerContract() mainchain.Addr {
	ctx := context.Background()
	channelAddrBundle := deploy.DeployAll(EtherBase.Auth, EtherBase.Client, ctx, 0)
	ledgerAddr := channelAddrBundle.CelerLedgerAddr

	// Disable channel deposit limit
	LogBlkNum(EtherBase.Client)
	ledgerContract, err := ledger.NewCelerLedger(ledgerAddr, EtherBase.Client)
	ChkErr(err, "failed to NewCelerLedger")
	tx, err := ledgerContract.DisableBalanceLimits(EtherBase.Auth)
	ChkErr(err, "failed disable channel deposit limits")
	WaitMinedWithChk(ctx, EtherBase.Client, tx, 0, "Disable balance limit")

	log.Infoln("Ledger address:", ledgerAddr.String())
	return ledgerAddr
}

func DeployERC20Contract() (mainchain.Addr, *mainchain.ERC20) {
	ctx := context.Background()
	initAmt := new(big.Int)
	initAmt.SetString("1"+strings.Repeat("0", 28), 10)
	erc20Addr, tx, erc20, err := mainchain.DeployERC20(EtherBase.Auth, EtherBase.Client, initAmt, "Celer", 18, "CELR")
	ChkErr(err, "failed to deploy ERC20")
	WaitMinedWithChk(ctx, EtherBase.Client, tx, 0, "Deploy ERC20 "+erc20Addr.Hex())

	log.Infoln("Erc20 address:", erc20Addr.String())
	return erc20Addr, erc20
}

func DeployGuardContract(sgnParams *SGNParams) mainchain.Addr {
	ctx := context.Background()
	guardAddr, tx, _, err := mainchain.DeployGuard(EtherBase.Auth, EtherBase.Client, sgnParams.CelrAddr, sgnParams.BlameTimeout, sgnParams.MinValidatorNum, sgnParams.MinStakingPool, sgnParams.SidechainGoLiveTimeout, sgnParams.MaxValidatorNum)
	ChkErr(err, "failed to deploy Guard contract")
	WaitMinedWithChk(ctx, EtherBase.Client, tx, 0, "Deploy Guard "+guardAddr.Hex())

	log.Infoln("Guard address:", guardAddr.String())
	return guardAddr
}

func DeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy contracts",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			viper.SetConfigFile(viper.GetString(common.FlagConfig))
			err = viper.ReadInConfig()
			if err != nil {
				return
			}

			ws := viper.GetString(common.FlagEthWS)
			err = DefaultEthClient.SetClient(ws)
			if err != nil {
				return
			}

			err = DefaultEthClient.SetAuth(viper.GetString(common.FlagEthKeystore), viper.GetString(common.FlagEthPassphrase))
			if err != nil {
				return
			}

			if ws == LocalGeth {
				SetEthBaseKs("./docker-volumes/geth-env")
				err = FundAddrsETH("1"+strings.Repeat("0", 20),
					[]mainchain.Addr{mainchain.Hex2Addr(ValEthAddrs[0]), mainchain.Hex2Addr(ValEthAddrs[1])})
				ChkErr(err, "fund client0 and client1")
			}

			ledgerAddr := DeployLedgerContract()
			viper.Set(common.FlagEthLedgerAddress, ledgerAddr)

			erc20Addr, erc20 := DeployERC20Contract()
			sgnParams := &SGNParams{
				BlameTimeout:           big.NewInt(50),
				MinValidatorNum:        big.NewInt(0),
				MinStakingPool:         big.NewInt(100),
				SidechainGoLiveTimeout: big.NewInt(0),
				CelrAddr:               erc20Addr,
				MaxValidatorNum:        big.NewInt(10),
			}
			guardAddr := DeployGuardContract(sgnParams)
			viper.Set(common.FlagEthGuardAddress, guardAddr)
			err = viper.WriteConfig()
			ChkErr(err, "failed to write config")

			if ws == LocalGeth {
				amt := new(big.Int)
				amt.SetString("1"+strings.Repeat("0", 19), 10)
				tx, err := erc20.Approve(DefaultEthClient.Auth, guardAddr, amt)
				ChkErr(err, "failed to approve erc20")
				WaitMinedWithChk(context.Background(), DefaultEthClient.Client, tx, 0, "approve erc20")
			}

			return nil
		},
	}

	return cmd
}
