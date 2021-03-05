package ops

import (
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/sgn"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func WithdrawCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Withdraw delegated stake",
	}
	cmd.AddCommand(
		intendWithdrawCommand(),
		confirmWithdrawCommand(),
		withdrawFromUnbondedCandidateCommand(),
		withdrawReward(cdc),
	)
	return cmd
}

func intendWithdrawCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intend",
		Short: "Send a withdrawal intent for the stake delegated to a candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}
			amount := calcRawAmount(viper.GetString(amountFlag))
			candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

			log.Infof(
				"Sending intent to withdraw amount %s from candidate %s",
				amount,
				candidate.Hex(),
			)
			_, err = ethClient.Transactor.TransactWaitMined(
				"IntendWithdraw",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.IntendWithdraw(opts, candidate, amount)
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.Flags().String(amountFlag, "", "Withdraw amount (integer in unit of CELR)")
	cmd.MarkFlagRequired(candidateFlag)
	cmd.MarkFlagRequired(amountFlag)
	return cmd
}

func confirmWithdrawCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm",
		Short: "Confirm withdrawal intents for the stake delegated to a candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}
			candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

			log.Infof(
				"Confirming withdrawal intents for the stake delegated to candidate %s",
				candidate.Hex(),
			)
			_, err = ethClient.Transactor.TransactWaitMined(
				"ConfirmWithdraw",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.ConfirmWithdraw(opts, candidate)
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.MarkFlagRequired(candidateFlag)
	return cmd
}

func withdrawFromUnbondedCandidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonded-candidate",
		Short: "Withdraw delegated stake from an unbonded candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}
			amount := calcRawAmount(viper.GetString(amountFlag))
			candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

			log.Infof(
				"Withdrawing amount %s delegated from an unbonded candidate %s",
				amount,
				candidate.Hex(),
			)
			_, err = ethClient.Transactor.TransactWaitMined(
				"WithdrawFromUnbondedCandidate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.WithdrawFromUnbondedCandidate(opts, candidate, amount)
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.Flags().String(amountFlag, "", "Withdraw amount (integer in unit of CELR)")
	cmd.MarkFlagRequired(candidateFlag)
	cmd.MarkFlagRequired(amountFlag)
	return cmd
}

func withdrawReward(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward [eth-address]",
		Short: "Withdraw reward on mainchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}

			cliCtx := common.NewQueryCLIContext(cdc)
			reward, err := validator.CLIQueryReward(cliCtx, validator.RouterKey, ethClient.Address.Hex())
			if err != nil {
				log.Errorln("query reward error", err)
				return err
			}

			if len(reward.RewardProtoBytes) == 0 {
				log.Info("no signed reward")
				return nil
			}

			var mingingReward, serviceReward *big.Int
			var pbReward sgn.Reward
			err = proto.Unmarshal(reward.RewardProtoBytes, &pbReward)
			if err != nil {
				log.Errorln("proto umarshal err", err, reward.RewardProtoBytes)
			} else {
				mingingReward = new(big.Int).SetBytes(pbReward.CumulativeMiningReward)
				serviceReward = new(big.Int).SetBytes(pbReward.CumulativeServiceReward)
			}

			var signers []mainchain.Addr
			for _, sigs := range reward.Sigs {
				signers = append(signers, mainchain.Hex2Addr(sigs.Signer))
			}
			signerStakes, _, quorumStakes, err := ethClient.CheckVotingPower(signers)
			if err != nil {
				log.Errorln("check signers voting power error", err)
				return err
			}
			if quorumStakes.Cmp(signerStakes) > 0 {
				log.Infof("signer stakes %s smaller than quorum stakes %s", signerStakes, quorumStakes)
				return nil
			}

			redeemedMiningReward, err := ethClient.DPoS.RedeemedMiningReward(&bind.CallOpts{}, ethClient.Address)
			if err != nil {
				log.Errorln("query RedeemedMiningReward err", err)
				return err
			}
			redeemedServiceReward, err := ethClient.SGN.RedeemedServiceReward(&bind.CallOpts{}, ethClient.Address)
			if err != nil {
				log.Errorln("query RedeemedServiceReward err", err)
				return err
			}
			if mingingReward.Cmp(redeemedMiningReward) <= 0 && serviceReward.Cmp(redeemedServiceReward) <= 0 {
				log.Info("no new reward")
				return nil
			}

			log.Infof("Withdrawing mining reward %s service reward %s",
				mingingReward.Sub(mingingReward, redeemedMiningReward), serviceReward.Sub(serviceReward, redeemedServiceReward))

			_, err = ethClient.Transactor.TransactWaitMined(
				"RedeemReward",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.SGN.RedeemReward(opts, reward.GetRewardRequest())
				},
			)
			if err != nil {
				return err
			}

			redeemedMiningReward, err = ethClient.DPoS.RedeemedMiningReward(&bind.CallOpts{}, ethClient.Address)
			if err != nil {
				log.Errorln("query RedeemedMiningReward err", err)
				return err
			}
			redeemedServiceReward, err = ethClient.SGN.RedeemedServiceReward(&bind.CallOpts{}, ethClient.Address)
			if err != nil {
				log.Errorln("query RedeemedServiceReward err", err)
				return err
			}
			log.Infof("Total withdrawn mining reward %s service reward %s", redeemedMiningReward, redeemedServiceReward)

			return nil
		},
	}

	return cmd
}
