package validator

import (
	"fmt"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	log.Infoln("Handle msg to initialize candidate", msg.EthAddress, "from sender", msg.Sender)

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	accAddress := sdk.AccAddress(candidateInfo.SidechainAddr)
	account := keeper.accountKeeper.GetAccount(ctx, accAddress)
	if account == nil {
		log.Infof("Set new account %x for candidate %s", accAddress, msg.EthAddress)
		account = keeper.accountKeeper.NewAccountWithAddress(ctx, accAddress)
		keeper.accountKeeper.SetAccount(ctx, account)
	}

	_, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		log.Infof("Created a new profile for candidate %s account %x", msg.EthAddress, accAddress)
		keeper.SetCandidate(ctx, msg.EthAddress, NewCandidate(accAddress))
	}

	return sdk.Result{}
}

// Handle a message to claim validator
func handleMsgClaimValidator(ctx sdk.Context, keeper Keeper, msg MsgClaimValidator) sdk.Result {
	log.Infof("Handle MsgClaimValidator. %+v", msg)

	pk, err := sdk.GetConsPubKeyBech32(msg.PubKey)
	if err != nil {
		return sdk.ErrInvalidPubKey(err.Error()).Result()
	}

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	if !mainchain.IsBonded(candidateInfo) {
		return sdk.ErrInternal("Candidate is not in validator set").Result()
	}

	if !sdk.AccAddress(candidateInfo.SidechainAddr).Equals(msg.Sender) {
		return sdk.ErrInternal("Sender has different address recorded on mainchain").Result()
	}

	candidate, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		return sdk.ErrInternal("Candidate does not exist").Result()
	}

	// Make sure both val address and pub address have not been used before
	valAddress := sdk.ValAddress(candidateInfo.SidechainAddr)
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
		keeper.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	}

	validator.Status = sdk.Bonded
	updateValidatorToken(ctx, keeper, validator, candidateInfo.StakingPool)

	candidate.Transactors = msg.Transactors
	keeper.SetCandidate(ctx, msg.EthAddress, candidate)

	return sdk.Result{}
}

// Handle a message to sync validator
func handleMsgSyncValidator(ctx sdk.Context, keeper Keeper, msg MsgSyncValidator) sdk.Result {
	log.Infof("Handle MsgSyncValidator. %+v", msg)

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	valAddress := sdk.ValAddress(candidateInfo.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return sdk.ErrInternal("Validator does not exist").Result()
	}

	updateValidatorToken(ctx, keeper, validator, candidateInfo.StakingPool)
	if !mainchain.IsBonded(candidateInfo) {
		validator.Status = mainchain.ParseStatus(candidateInfo)
		keeper.stakingKeeper.SetValidator(ctx, validator)
		keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	}

	return sdk.Result{}
}

// Handle a message to sync delegator
func handleMsgSyncDelegator(ctx sdk.Context, keeper Keeper, msg MsgSyncDelegator) sdk.Result {
	log.Infof("Handle MsgSyncDelegator. %+v", msg)
	delegator := keeper.GetDelegator(ctx, msg.CandidateAddress, msg.DelegatorAddress)
	di, err := keeper.ethClient.Guard.GetDelegatorInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, mainchain.Hex2Addr(msg.CandidateAddress), mainchain.Hex2Addr(msg.DelegatorAddress))
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
	log.Infof("Handle MsgWithdrawReward. %+v", msg)
	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return sdk.ErrInternal("Reward does not exist").Result()
	}
	if !reward.HasNewReward() {
		return sdk.ErrInternal("No new reward").Result()
	}

	reward.InitateWithdraw()
	keeper.SetReward(ctx, msg.EthAddress, reward)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			ModuleName,
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
	log.Infof("Handle MsgSignReward. EthAddress %s, Sender %s", msg.EthAddress, msg.Sender.String())
	validator, found := keeper.stakingKeeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return sdk.ErrInternal("Sender is not validator").Result()
	}
	if validator.Status != sdk.Bonded {
		return sdk.ErrInternal("Validator is not bonded").Result()
	}

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return sdk.ErrInternal("Reward does not exist").Result()
	}

	err := reward.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to add sig: %s", err)).Result()
	}

	keeper.SetReward(ctx, msg.EthAddress, reward)
	return sdk.Result{}
}

func updateValidatorToken(ctx sdk.Context, keeper Keeper, validator staking.Validator, totalTokens *big.Int) {
	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = sdk.NewIntFromBigInt(totalTokens)
	validator.DelegatorShares = validator.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)
	keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
}
