package validator

import (
	"fmt"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// NewHandler returns a handler for "validator" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		var res *sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgUpdateSidechainAddr:
			res, err = handleMsgUpdateSidechainAddr(ctx, keeper, msg, logEntry)
		case MsgSetTransactors:
			res, err = handleMsgSetTransactors(ctx, keeper, msg, logEntry)
		case MsgClaimValidator:
			res, err = handleMsgClaimValidator(ctx, keeper, msg, logEntry)
		case MsgSyncValidator:
			res, err = handleMsgSyncValidator(ctx, keeper, msg, logEntry)
		case MsgSyncDelegator:
			res, err = handleMsgSyncDelegator(ctx, keeper, msg, logEntry)
		case MsgWithdrawReward:
			res, err = handleMsgWithdrawReward(ctx, keeper, msg, logEntry)
		case MsgSignReward:
			res, err = handleMsgSignReward(ctx, keeper, msg, logEntry)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}

		if err != nil {
			logEntry.Error = append(logEntry.Error, err.Error())
		}

		seal.CommitMsgLog(logEntry)
		return res, err
	}
}

// Handle a message to update sidechain address
func handleMsgUpdateSidechainAddr(ctx sdk.Context, keeper Keeper, msg MsgUpdateSidechainAddr, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to query candidate profile: %s", err)
	}

	if !candidateInfo.DPoSCandidateInfo.Initialized {
		return nil, fmt.Errorf("Candidate has not been initialized")
	}

	accAddress := sdk.AccAddress(candidateInfo.SidechainAddr)
	InitAccount(ctx, keeper, accAddress)

	// TODO: only handle the case of first update (initialization), need to handle the case of replacing sidechain address
	_, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		log.Infof("Created a new profile for candidate %s account %x", msg.EthAddress, accAddress)
		keeper.SetCandidate(ctx, NewCandidate(msg.EthAddress, accAddress))
	}

	return &sdk.Result{}, nil
}

// Handle a message to set transactors
func handleMsgSetTransactors(ctx sdk.Context, keeper Keeper, msg MsgSetTransactors, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	for _, transactor := range msg.Transactors {
		logEntry.Transactor = append(logEntry.Transactor, transactor.String())
	}

	candidate, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		return nil, fmt.Errorf("Candidate does not exist")
	}

	if !candidate.Operator.Equals(msg.Sender) {
		return nil, fmt.Errorf("The candidate is not operated by the sender.")
	}

	candidate.Transactors = msg.Transactors
	for _, transactor := range candidate.Transactors {
		InitAccount(ctx, keeper, transactor)
	}

	keeper.SetCandidate(ctx, candidate)
	return &sdk.Result{}, nil
}

// Handle a message to claim validator
func handleMsgClaimValidator(ctx sdk.Context, keeper Keeper, msg MsgClaimValidator, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress
	logEntry.PubKey = msg.PubKey

	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, msg.PubKey)
	if err != nil {
		return nil, fmt.Errorf("GetConsPubKeyBech32 err: %s", err)
	}

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to query candidate profile: %s", err)
	}

	if !mainchain.IsBonded(candidateInfo.DPoSCandidateInfo) {
		return nil, fmt.Errorf("Candidate is not in validator set")
	}

	if !sdk.AccAddress(candidateInfo.SidechainAddr).Equals(msg.Sender) {
		return nil, fmt.Errorf("Sender has different address recorded on mainchain. mainchain record: %x; sender: %x", sdk.AccAddress(candidateInfo.SidechainAddr), msg.Sender)
	}

	// Make sure both val address and pub address have not been used before
	valAddress := sdk.ValAddress(candidateInfo.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	_, f := keeper.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk))
	if found != f {
		return nil, fmt.Errorf("Invalid sender address or public key")
	}

	if !found {
		description := staking.Description{
			Moniker:  msg.EthAddress,
			Identity: msg.EthAddress,
		}
		validator = staking.NewValidator(valAddress, pk, description)
		keeper.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	}

	updateValidatorToken(ctx, keeper, validator, candidateInfo)
	return &sdk.Result{}, nil
}

// Handle a message to sync validator
func handleMsgSyncValidator(ctx sdk.Context, keeper Keeper, msg MsgSyncValidator, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to query candidate profile: %s", err)
	}

	valAddress := sdk.ValAddress(candidateInfo.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return &sdk.Result{}, fmt.Errorf("Validator does not exist")
	}

	updateValidatorToken(ctx, keeper, validator, candidateInfo)
	return &sdk.Result{}, nil
}

// Handle a message to sync delegator
func handleMsgSyncDelegator(ctx sdk.Context, keeper Keeper, msg MsgSyncDelegator, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.CandidateAddr = msg.CandidateAddress
	logEntry.DelegatorAddr = msg.DelegatorAddress

	delegator := keeper.GetDelegator(ctx, msg.CandidateAddress, msg.DelegatorAddress)
	di, err := keeper.ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, mainchain.Hex2Addr(msg.CandidateAddress), mainchain.Hex2Addr(msg.DelegatorAddress))
	if err != nil {
		return nil, fmt.Errorf("Failed to query delegator info: %s", err)
	}

	delegator.DelegatedStake = sdk.NewIntFromBigInt(di.DelegatedStake)
	keeper.SetDelegator(ctx, msg.CandidateAddress, msg.DelegatorAddress, delegator)
	keeper.SnapshotCandidate(ctx, msg.CandidateAddress)

	return &sdk.Result{}, nil
}

// Handle a message to withdraw reward
func handleMsgWithdrawReward(ctx sdk.Context, keeper Keeper, msg MsgWithdrawReward, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return nil, fmt.Errorf("Reward does not exist")
	}
	if !reward.HasNewReward() {
		return nil, fmt.Errorf("No new reward")
	}

	reward.InitateWithdraw()
	keeper.SetReward(ctx, reward)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			ModuleName,
			sdk.NewAttribute(sdk.AttributeKeyAction, ActionInitiateWithdraw),
			sdk.NewAttribute(AttributeKeyEthAddress, msg.EthAddress),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// Handle a message to sign reward
func handleMsgSignReward(ctx sdk.Context, keeper Keeper, msg MsgSignReward, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	validator, found := keeper.stakingKeeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return nil, fmt.Errorf("Sender is not validator")
	}
	if validator.Status != sdk.Bonded {
		return nil, fmt.Errorf("Validator is not bonded")
	}

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return nil, fmt.Errorf("Reward does not exist")
	}

	err := reward.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return nil, fmt.Errorf("Failed to add sig: %s", err)
	}

	keeper.SetReward(ctx, reward)
	return &sdk.Result{}, nil
}

func updateValidatorToken(ctx sdk.Context, keeper Keeper, validator staking.Validator, candidateInfo mainchain.CandidateInfo) {
	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = sdk.NewIntFromBigInt(candidateInfo.DPoSCandidateInfo.StakingPool).QuoRaw(common.TokenDec)
	validator.Status = mainchain.ParseStatus(candidateInfo.DPoSCandidateInfo)
	validator.DelegatorShares = validator.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)

	if validator.Status == sdk.Bonded {
		keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
	}
}
