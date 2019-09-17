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
		default:
			errMsg := fmt.Sprintf("Unrecognized validator Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to initialize candidate
func handleMsgInitializeCandidate(ctx sdk.Context, keeper Keeper, msg MsgInitializeCandidate) sdk.Result {
	cp, err := keeper.ethClient.Guard.GetCandidateInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(msg.EthAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	accAddress := sdk.AccAddress(cp.SidechainAddr)
	account := keeper.accountKeeper.NewAccountWithAddress(ctx, accAddress)
	keeper.accountKeeper.SetAccount(ctx, account)
	return sdk.Result{}
}

// Handle a message to claim validator
func handleMsgClaimValidator(ctx sdk.Context, keeper Keeper, msg MsgClaimValidator) sdk.Result {
	pk, err := sdk.GetConsPubKeyBech32(msg.PubKey)
	if err != nil {
		return sdk.ErrInvalidPubKey(err.Error()).Result()
	}

	cp, err := keeper.ethClient.Guard.GetCandidateInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(msg.EthAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	if !cp.IsVldt {
		return sdk.ErrInternal("Candidate is not in validator set").Result()
	}

	if !sdk.AccAddress(cp.SidechainAddr).Equals(msg.Sender) {
		return sdk.ErrInternal("Sender has different address recorded on mainchain").Result()
	}

	valAddress := sdk.ValAddress(cp.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	_, f := keeper.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk))

	if found != f {
		return sdk.ErrInternal("Invalid sender address or public key").Result()
	}

	if !found {
		description := staking.Description{
			Moniker: msg.EthAddress,
		}
		validator = staking.NewValidator(valAddress, pk, description)
		validator.Status = sdk.Bonded
		keeper.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	}

	updateValidatorToken(ctx, keeper, validator, cp.TotalStake)
	return sdk.Result{}
}

// Handle a message to sync validator
func handleMsgSyncValidator(ctx sdk.Context, keeper Keeper, msg MsgSyncValidator) sdk.Result {
	cp, err := keeper.ethClient.Guard.GetCandidateInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(msg.EthAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	valAddress := sdk.ValAddress(cp.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return sdk.ErrInternal("Validator does not exist").Result()
	}

	if !cp.IsVldt {
		keeper.stakingKeeper.RemoveValidator(ctx, validator.OperatorAddress)
	}

	updateValidatorToken(ctx, keeper, validator, cp.TotalStake)
	return sdk.Result{}
}

func updateValidatorToken(ctx sdk.Context, keeper Keeper, validator staking.Validator, totalTokens *big.Int) {
	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = sdk.NewIntFromBigInt(totalTokens)
	validator.DelegatorShares = validator.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)
	keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
}
