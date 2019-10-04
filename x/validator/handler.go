package validator

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// NewHandler returns a handler for "validator" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgInitializeCandidate:
			return handleMsgInitializeCandidate(ctx, keeper, msg)
		case MsgClaimValidator:
			return handleMsgClaimValidator(ctx, keeper, msg)
		case MsgSyncValidator:
			return handleMsgSyncValidator(ctx, keeper, msg)
		case MsgSyncDelegator:
			return handleMsgSyncDelegator(ctx, keeper, msg)
		case MsgWithdrawReward:
			return handleMsgWithdrawReward(ctx, keeper, msg)
		case MsgSignReward:
			return handleMsgSignReward(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized validator Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to initialize candidate
func handleMsgInitializeCandidate(ctx sdk.Context, keeper Keeper, msg MsgInitializeCandidate) sdk.Result {
	cp, err := GetCandidateInfo(ctx, keeper, msg.EthAddress)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	accAddress := sdk.AccAddress(cp.SidechainAddr)
	account := keeper.accountKeeper.GetAccount(ctx, accAddress)
	if account == nil {
		account = keeper.accountKeeper.NewAccountWithAddress(ctx, accAddress)
		keeper.accountKeeper.SetAccount(ctx, account)
	}

	_, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		keeper.SetCandidate(ctx, msg.EthAddress, NewCandidate())
	}
	return sdk.Result{}
}

// Handle a message to claim validator
func handleMsgClaimValidator(ctx sdk.Context, keeper Keeper, msg MsgClaimValidator) sdk.Result {
	pk, err := sdk.GetConsPubKeyBech32(msg.PubKey)
	if err != nil {
		return sdk.ErrInvalidPubKey(err.Error()).Result()
	}

	cp, err := GetCandidateInfo(ctx, keeper, msg.EthAddress)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	if !cp.IsVldt {
		return sdk.ErrInternal("Candidate is not in validator set").Result()
	}

	if !sdk.AccAddress(cp.SidechainAddr).Equals(msg.Sender) {
		return sdk.ErrInternal("Sender has different address recorded on mainchain").Result()
	}

	// Make sure both val address and pub address have not been used before
	valAddress := sdk.ValAddress(cp.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	_, f := keeper.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk))

	if found != f {
		return sdk.ErrInternal("Invalid sender address or public key").Result()
	}

	if !found {
		description := staking.Description{
			Moniker:  msg.EthAddress,
			Identity: msg.EthAddress,
		}
		validator = staking.NewValidator(valAddress, pk, description)
		validator.Status = sdk.Bonded
		keeper.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	}

	updateValidatorToken(ctx, keeper, validator, cp.StakingPool)
	return sdk.Result{}
}

// Handle a message to sync validator
func handleMsgSyncValidator(ctx sdk.Context, keeper Keeper, msg MsgSyncValidator) sdk.Result {
	cp, err := GetCandidateInfo(ctx, keeper, msg.EthAddress)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	valAddress := sdk.ValAddress(cp.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return sdk.ErrInternal("Validator does not exist").Result()
	}

	if cp.IsVldt {
		updateValidatorToken(ctx, keeper, validator, cp.StakingPool)
	} else {
		keeper.stakingKeeper.RemoveValidator(ctx, validator.OperatorAddress)
	}

	return sdk.Result{}
}

// Handle a message to sync delegator
func handleMsgSyncDelegator(ctx sdk.Context, keeper Keeper, msg MsgSyncDelegator) sdk.Result {
	delegator := keeper.GetDelegator(ctx, msg.CandidateAddress, msg.DelegatorAddress)

	di, err := keeper.ethClient.Guard.GetDelegatorInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(msg.CandidateAddress), ethcommon.HexToAddress(msg.DelegatorAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query delegator info: %s", err)).Result()
	}

	delegator.DelegatedStake = sdk.NewIntFromBigInt(di.DelegatedStake)
	keeper.SetDelegator(ctx, msg.CandidateAddress, msg.DelegatorAddress, delegator)
	keeper.SnapshotCandidate(ctx, msg.CandidateAddress)

	return sdk.Result{}
}

// Handle a message to withdraw reward
func handleMsgWithdrawReward(ctx sdk.Context, keeper Keeper, msg MsgWithdrawReward) sdk.Result {
	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return sdk.ErrInternal("Reward does not exist").Result()
	}
	if !reward.HasNewReward() {
		return sdk.ErrInternal("No new reward").Result()
	}

	reward.InitateWithdraw()
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, ActionInitiateWithdraw),
			sdk.NewAttribute(AttributeKeyEthAddress, msg.EthAddress),
		),
	)
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Handle a message to sign reward
func handleMsgSignReward(ctx sdk.Context, keeper Keeper, msg MsgSignReward) sdk.Result {
	validator, found := keeper.stakingKeeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return sdk.ErrInternal("Sender is not validator").Result()
	}

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return sdk.ErrInternal("Reward does not exist").Result()
	}

	err := reward.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to add sig: %s", err)).Result()
	}

	return sdk.Result{}
}

func updateValidatorToken(ctx sdk.Context, keeper Keeper, validator staking.Validator, totalTokens *big.Int) {
	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = sdk.NewIntFromBigInt(totalTokens)
	validator.DelegatorShares = validator.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)
	keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
}
