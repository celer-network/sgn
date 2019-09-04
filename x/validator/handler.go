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
		case MsgSyncValidator:
			return handleMsgSyncValidator(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized validator Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set eth address
func handleMsgSyncValidator(ctx sdk.Context, keeper Keeper, msg MsgSyncValidator) sdk.Result {
	pk, err := sdk.GetConsPubKeyBech32(msg.PubKey)
	if err != nil {
		return sdk.ErrInvalidPubKey(err.Error()).Result()
	}

	cp, err := keeper.ethClient.Guard.CandidateProfiles(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(msg.EthAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	if !sdk.AccAddress(cp.SidechainAddr).Equals(msg.Sender) {
		return sdk.ErrInternal("Sender is not selected validator").Result()
	}

	valAddress := sdk.ValAddress(msg.Sender)
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
		keeper.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	}

	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = sdk.NewIntFromBigInt(cp.Stakes)
	validator.DelegatorShares = validator.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)
	keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)

	return sdk.Result{}
}
