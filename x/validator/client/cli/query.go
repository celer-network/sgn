package cli

import (
	"fmt"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/x/validator/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

const (
	flagSeq = "seq"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	validatorQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the validator module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	validatorQueryCmd.AddCommand(flags.GetCommands(
		GetCmdCandidate(storeKey, cdc),
		GetCmdDelegator(storeKey, cdc),
		GetCmdValidator(staking.StoreKey, cdc),
		GetCmdValidators(staking.StoreKey, cdc),
		GetCmdSyncer(storeKey, cdc),
		GetCmdReward(storeKey, cdc),
		GetCmdRewardRequest(storeKey, cdc),
	)...)
	return validatorQueryCmd
}

// GetCmdSyncer queries syncer info
func GetCmdSyncer(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "syncer",
		Short: "query syncer info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			syncer, err := QuerySyncer(cliCtx, queryRoute)
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(syncer)
		},
	}
}

// Query syncer info
func QuerySyncer(cliCtx context.CLIContext, queryRoute string) (syncer types.Syncer, err error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySyncer)
	res, err := common.RobustQuery(cliCtx, route)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &syncer)
	return
}

// GetCmdDelegator queries request info
func GetCmdDelegator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delegator [candidateEthAddr] [delegatorEthAddr]",
		Short: "query delegator info by candidateEthAddr and delegatorEthAddr",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			delegator, err := QueryDelegator(cliCtx, queryRoute, args[0], args[1])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(delegator)
		},
	}
}

func QueryDelegator(cliCtx context.CLIContext, queryRoute, candidateAddress, delegatorAddress string) (delegator types.Delegator, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQueryDelegatorParams(candidateAddress, delegatorAddress))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegator)
	res, err := common.RobustQueryWithData(cliCtx, route, data)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &delegator)
	return
}

// GetCmdCandidate queries request info
func GetCmdCandidate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "candidate [candidateEthAddr]",
		Short: "query candidate info by candidateEthAddr",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			candidate, err := QueryCandidate(cliCtx, queryRoute, args[0])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(candidate)
		},
	}
}

func QueryCandidate(cliCtx context.CLIContext, queryRoute, ethAddress string) (candidate types.Candidate, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQueryCandidateParams(ethAddress))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryCandidate)
	res, err := common.RobustQueryWithData(cliCtx, route, data)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &candidate)
	return
}

// GetCmdValidator queries validator info
func GetCmdValidator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "validator [validator account address]",
		Short: "query a validator by account address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			validator, err := QueryValidator(cliCtx, queryRoute, args[0])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}
			output := getValidatorOutput(&validator)
			return cliCtx.PrintOutput(output)
		},
	}
}

// GetCmdValidator queries validator info
func GetCmdValidators(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "validators",
		Short: "query all validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			validators, err := QueryValidators(cliCtx, queryRoute)
			if err != nil {
				log.Errorln("query error", err)
				return err
			}
			var outputs []ValidatorOutput
			for _, v := range validators {
				outputs = append(outputs, getValidatorOutput(&v))
			}
			return cliCtx.PrintOutput(outputs)
		},
	}
}

// QueryValidators is an interface for convenience to query (all) validators in staking module
func QueryValidators(cliCtx context.CLIContext, storeName string) (validators stakingTypes.Validators, err error) {
	resKVs, _, err := cliCtx.QuerySubspace(stakingTypes.ValidatorsKey, storeName)
	if err != nil {
		return
	}

	for _, kv := range resKVs {
		validators = append(validators, stakingTypes.MustUnmarshalValidator(cliCtx.Codec, kv.Value))
	}
	return
}

// QueryBondedValidators is an interface for convenience to query bonded validators in staking module
func QueryBondedValidators(cliCtx context.CLIContext, storeName string) (validators stakingTypes.Validators, err error) {
	allValidators, err := QueryValidators(cliCtx, storeName)
	if err != nil {
		return
	}

	for _, val := range allValidators {
		if val.Status == sdk.Bonded {
			validators = append(validators, val)
		}
	}

	return
}

// addrStr should be bech32 cosmos account address with prefix cosmos
func QueryValidator(cliCtx context.CLIContext, storeName string, addrStr string) (validator stakingTypes.Validator, err error) {
	addr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		log.Error(err)
		return
	}

	res, _, err := cliCtx.QueryStore(stakingTypes.GetValidatorKey(sdk.ValAddress(addr)), storeName)
	if err != nil {
		return
	}

	if len(res) == 0 {
		err = fmt.Errorf("%w for address %s", common.ErrRecordNotFound, addr)
		return
	}

	validator = stakingTypes.MustUnmarshalValidator(cliCtx.Codec, res)
	return
}

// GetCmdReward queries reward info
func GetCmdReward(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reward [ethAddress]",
		Short: "query reward info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			reward, err := QueryReward(cliCtx, queryRoute, args[0])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(reward)
		},
	}
}

// GetCmdRewardRequest queries reward request
func GetCmdRewardRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reward-request [ethAddress]",
		Short: "query reward request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			reward, err := QueryReward(cliCtx, queryRoute, args[0])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			log.Info(string(reward.GetRewardRequest()))
			return nil
		},
	}
}

// Query reward info
func QueryReward(cliCtx context.CLIContext, queryRoute string, ethAddress string) (reward types.Reward, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQueryRewardParams(ethAddress))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryReward)
	res, err := common.RobustQueryWithData(cliCtx, route, data)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &reward)
	return
}

type ValidatorOutput struct {
	AccountAddress  sdk.AccAddress      `json:"account_address" yaml:"account_address"`   // address of the validator's account; bech encoded in JSON
	OperatorAddress sdk.ValAddress      `json:"operator_address" yaml:"operator_address"` // address of the validator's operator; bech encoded in JSON
	ConsPubKey      string              `json:"consensus_pubkey" yaml:"consensus_pubkey"` // the consensus public key of the validator; bech encoded in JSON
	Status          string              `json:"status" yaml:"status"`                     // validator status (bonded/unbonding/unbonded)
	Tokens          sdk.Int             `json:"tokens" yaml:"tokens"`                     // delegated tokens (incl. self-delegation)
	DelegatorShares sdk.Dec             `json:"delegator_shares" yaml:"delegator_shares"` // total shares issued to a validator's delegators
	CommissionRate  sdk.Dec             `json:"commissionRate"`
	Description     staking.Description `json:"description" yaml:"description"` // description terms for the validator
}

func getValidatorOutput(v *stakingTypes.Validator) ValidatorOutput {
	bechConsPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, v.ConsPubKey)
	if err != nil {
		log.Error(err)
	}
	status := sdk.BondStatusUnbonded
	if v.Status == sdk.Unbonding {
		status = sdk.BondStatusUnbonding
	} else if v.Status == sdk.Bonded {
		status = sdk.BondStatusBonded
	}
	return ValidatorOutput{
		AccountAddress:  sdk.AccAddress(v.OperatorAddress),
		OperatorAddress: v.OperatorAddress,
		ConsPubKey:      bechConsPubKey,
		Status:          status,
		Tokens:          v.Tokens,
		DelegatorShares: v.DelegatorShares,
		CommissionRate:  v.Commission.CommissionRates.Rate,
		Description:     v.Description,
	}
}
