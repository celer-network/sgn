package validator

import (
	"fmt"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// NewHandler returns a handler for "validator" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		logEntry := seal.NewMsgLog()
		var res sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgInitializeCandidate:
			res, err = handleMsgInitializeCandidate(ctx, keeper, msg, logEntry)
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
			errMsg := fmt.Sprintf("Unrecognized validator Msg type: %v", msg.Type())
			log.Error(errMsg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}

		if err != nil {
			logEntry.Error = append(logEntry.Error, err.Error())
			seal.CommitMsgLog(logEntry)
			return sdk.ErrInternal(err.Error()).Result()
		}
		seal.CommitMsgLog(logEntry)
		return res
	}
}

// Handle a message to initialize candidate
func handleMsgInitializeCandidate(ctx sdk.Context, keeper Keeper, msg MsgInitializeCandidate, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	res := sdk.Result{}
	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return res, fmt.Errorf("Failed to query candidate profile: %s", err)
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

	return res, nil
}

// Handle a message to claim validator
func handleMsgClaimValidator(ctx sdk.Context, keeper Keeper, msg MsgClaimValidator, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress
	logEntry.PubKey = msg.PubKey
	for _, transactor := range msg.Transactors {
		logEntry.Transactor = append(logEntry.Transactor, transactor.String())
	}

	res := sdk.Result{}
	pk, err := sdk.GetConsPubKeyBech32(msg.PubKey)
	if err != nil {
		return res, fmt.Errorf("GetConsPubKeyBech32 err: %s", err)
	}

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return res, fmt.Errorf("Failed to query candidate profile: %s", err)
	}

	if !mainchain.IsBonded(candidateInfo) {
		return res, fmt.Errorf("Candidate is not in validator set")
	}

	if !sdk.AccAddress(candidateInfo.SidechainAddr).Equals(msg.Sender) {
		return res, fmt.Errorf("Sender has different address recorded on mainchain. mainchain record: %x; sender: %x", sdk.AccAddress(candidateInfo.SidechainAddr), msg.Sender)
	}

	candidate, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		return res, fmt.Errorf("Candidate does not exist")
	}

	// Make sure both val address and pub address have not been used before
	valAddress := sdk.ValAddress(candidateInfo.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	_, f := keeper.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk))
	if found != f {
		return res, fmt.Errorf("Invalid sender address or public key")
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

	return res, nil
}

// Handle a message to sync validator
func handleMsgSyncValidator(ctx sdk.Context, keeper Keeper, msg MsgSyncValidator, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	res := sdk.Result{}
	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return res, fmt.Errorf("Failed to query candidate profile: %s", err)
	}

	valAddress := sdk.ValAddress(candidateInfo.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return res, fmt.Errorf("Validator does not exist")
	}

	updateValidatorToken(ctx, keeper, validator, candidateInfo.StakingPool)
	if !mainchain.IsBonded(candidateInfo) {
		validator.Status = mainchain.ParseStatus(candidateInfo)
		keeper.stakingKeeper.SetValidator(ctx, validator)
		keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	}

	return res, nil
}

// Handle a message to sync delegator
func handleMsgSyncDelegator(ctx sdk.Context, keeper Keeper, msg MsgSyncDelegator, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.CandidateAddr = msg.CandidateAddress
	logEntry.DelegatorAddr = msg.DelegatorAddress

	res := sdk.Result{}
	delegator := keeper.GetDelegator(ctx, msg.CandidateAddress, msg.DelegatorAddress)
	di, err := keeper.ethClient.Guard.GetDelegatorInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, mainchain.Hex2Addr(msg.CandidateAddress), mainchain.Hex2Addr(msg.DelegatorAddress))
	if err != nil {
		return res, fmt.Errorf("Failed to query delegator info: %s", err)
	}

	delegator.DelegatedStake = sdk.NewIntFromBigInt(di.DelegatedStake)
	keeper.SetDelegator(ctx, msg.CandidateAddress, msg.DelegatorAddress, delegator)
	keeper.SnapshotCandidate(ctx, msg.CandidateAddress)

	return res, nil
}

// Handle a message to withdraw reward
func handleMsgWithdrawReward(ctx sdk.Context, keeper Keeper, msg MsgWithdrawReward, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	res := sdk.Result{}
	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return res, fmt.Errorf("Reward does not exist")
	}
	if !reward.HasNewReward() {
		return res, fmt.Errorf("No new reward")
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
	}, nil
}

// Handle a message to sign reward
func handleMsgSignReward(ctx sdk.Context, keeper Keeper, msg MsgSignReward, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	res := sdk.Result{}
	validator, found := keeper.stakingKeeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return res, fmt.Errorf("Sender is not validator")
	}
	if validator.Status != sdk.Bonded {
		return res, fmt.Errorf("Validator is not bonded")
	}

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return res, fmt.Errorf("Reward does not exist")
	}

	err := reward.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return res, fmt.Errorf("Failed to add sig: %s", err)
	}

	keeper.SetReward(ctx, msg.EthAddress, reward)
	return res, nil
}

func updateValidatorToken(ctx sdk.Context, keeper Keeper, validator staking.Validator, totalTokens *big.Int) {
	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = sdk.NewIntFromBigInt(totalTokens)
	validator.DelegatorShares = validator.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)
	keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
}
