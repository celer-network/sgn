package common

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn-contract/bindings/go/sgncontracts"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/testing/channel-eth-go/deploy"
	"github.com/celer-network/sgn/testing/channel-eth-go/ledger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeployLedgerContract() mainchain.Addr {
	ctx := context.Background()
	channelAddrBundle := deploy.DeployAll(EtherBaseAuth, EthClient, ctx, 0)
	ledgerAddr := channelAddrBundle.CelerLedgerAddr

	// Disable channel deposit limit
	LogBlkNum(EthClient)
	ledgerContract, err := ledger.NewCelerLedger(ledgerAddr, EthClient)
	ChkErr(err, "failed to NewCelerLedger")
	_, err = ledgerContract.DisableBalanceLimits(EtherBaseAuth)
	ChkErr(err, "failed disable channel deposit limits")

	log.Infoln("Ledger address:", ledgerAddr.String())
	return ledgerAddr
}

func DeployERC20Contract() (*types.Transaction, mainchain.Addr, *mainchain.ERC20) {
	initAmt := new(big.Int)
	initAmt.SetString("1"+strings.Repeat("0", 28), 10)
	erc20Addr, tx, erc20, err := mainchain.DeployERC20(EtherBaseAuth, EthClient, initAmt, "Celer", 18, "CELR")
	ChkErr(err, "failed to deploy ERC20")

	log.Infoln("Erc20 address:", erc20Addr.String())
	return tx, erc20Addr, erc20
}

func DeployDPoSSGNContracts(sgnParams *SGNParams) (*types.Transaction, mainchain.Addr, mainchain.Addr) {
	dposAddr, _, _, err := sgncontracts.DeployDPoS(
		EtherBaseAuth,
		EthClient,
		sgnParams.CelrAddr,
		sgnParams.GovernProposalDeposit,
		sgnParams.GovernVoteTimeout,
		sgnParams.SlashTimeout,
		sgnParams.MinValidatorNum,
		sgnParams.MaxValidatorNum,
		sgnParams.MinStakingPool,
		sgnParams.AdvanceNoticePeriod,
		sgnParams.SidechainGoLiveTimeout)
	ChkErr(err, "failed to deploy DPoS contract")

	sgnAddr, _, _, err := sgncontracts.DeploySGN(EtherBaseAuth, EthClient, sgnParams.CelrAddr, dposAddr)
	ChkErr(err, "failed to deploy SGN contract")

	// TODO: register SGN address on DPoS contract
	dpos, err := sgncontracts.NewDPoS(dposAddr, EthClient)
	ChkErr(err, "failed to new DPoS instance")
	EtherBaseAuth.GasLimit = 8000000
	tx, err := dpos.RegisterSidechain(EtherBaseAuth, sgnAddr)
	EtherBaseAuth.GasLimit = 0
	ChkErr(err, "failed to register SGN address on DPoS contract")

	log.Infoln("DPoS address:", dposAddr.String())
	log.Infoln("SGN address:", sgnAddr.String())

	return tx, dposAddr, sgnAddr
}

func DeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy contracts",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			configFileViper := viper.New()
			configFileViper.SetConfigFile(viper.GetString(common.FlagConfig))
			err = configFileViper.ReadInConfig()
			if err != nil {
				return err
			}
			ethurl := configFileViper.GetString(common.FlagEthGateway)
			var rpcClient *rpc.Client
			rpcClient, err = rpc.Dial(ethurl)
			if err != nil {
				return err
			}
			EthClient = ethclient.NewClient(rpcClient)

			var ksBytes []byte
			ksBytes, err = ioutil.ReadFile(configFileViper.GetString(common.FlagEthKeystore))
			if err != nil {
				return err
			}
			EtherBaseAuth, err = bind.NewTransactor(
				strings.NewReader(string(ksBytes)), configFileViper.GetString(common.FlagEthPassphrase))
			if err != nil {
				return err
			}

			if ethurl == LocalGeth {
				SetEthBaseKs("./docker-volumes/geth-env")
				err = FundAddrsETH("1"+strings.Repeat("0", 20),
					[]mainchain.Addr{
						mainchain.Hex2Addr(ValEthAddrs[0]),
						mainchain.Hex2Addr(ClientEthAddrs[0]),
						mainchain.Hex2Addr(ClientEthAddrs[1]),
					})
				ChkErr(err, "fund ETH to validator and clients")
			}

			ledgerAddr := DeployLedgerContract()

			if ethurl == LocalGeth {
				genesisPath := os.ExpandEnv("$HOME/.sgnd/config/genesis.json")
				genesisViper := viper.New()
				genesisViper.SetConfigFile(genesisPath)
				err := genesisViper.ReadInConfig()
				ChkErr(err, "Failed to read genesis")
				genesisViper.Set("app_state.guard.params.ledger_address", ledgerAddr.Hex())
				err = genesisViper.WriteConfig()
				ChkErr(err, "Failed to write genesis")
			}

			_, erc20Addr, erc20 := DeployERC20Contract()
			// NOTE: values below are for local tests
			sgnParams := &SGNParams{
				CelrAddr:               erc20Addr,
				GovernProposalDeposit:  big.NewInt(1000000000000000000), // 1 CELR
				GovernVoteTimeout:      big.NewInt(90),
				SlashTimeout:           big.NewInt(15),
				MinValidatorNum:        big.NewInt(1),
				MaxValidatorNum:        big.NewInt(5),
				MinStakingPool:         big.NewInt(5000000000000000000), // 5 CELR
				AdvanceNoticePeriod:    big.NewInt(30),
				SidechainGoLiveTimeout: big.NewInt(0),
			}
			tx, dposAddr, sgnAddr := DeployDPoSSGNContracts(sgnParams)
			WaitMinedWithChk(context.Background(), EthClient, tx, BlockDelay, PollingInterval, "DeployDPoSContracts")

			configFileViper.Set(common.FlagEthCelrAddress, erc20Addr.Hex())
			configFileViper.Set(common.FlagEthDPoSAddress, dposAddr.Hex())
			configFileViper.Set(common.FlagEthSGNAddress, sgnAddr.Hex())
			err = configFileViper.WriteConfig()
			ChkErr(err, "failed to write config")

			if ethurl == LocalGeth {
				amt := new(big.Int)
				amt.SetString("1"+strings.Repeat("0", 19), 10)
				tx, err := erc20.Approve(EtherBaseAuth, dposAddr, amt)
				ChkErr(err, "failed to approve erc20")
				WaitMinedWithChk(context.Background(), EthClient, tx, BlockDelay, PollingInterval, "approve erc20")
				err = FundAddrsErc20(erc20Addr,
					[]mainchain.Addr{
						mainchain.Hex2Addr(ClientEthAddrs[0]),
						mainchain.Hex2Addr(ClientEthAddrs[1]),
					},
					"1"+strings.Repeat("0", 20),
				)
				ChkErr(err, "fund test CELR to clients")
			}

			return nil
		},
	}

	return cmd
}
