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

var (
	okRes  = sdk.Result{}
	errRes = sdk.Result{Code: sdk.CodeInternal}
)

var PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil))

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
			if res.IsOK() {
				logEntry.Warn = append(logEntry.Warn, err.Error())
			} else {
				logEntry.Error = append(logEntry.Error, err.Error())
			}

		}
		seal.CommitMsgLog(logEntry)

		if !res.IsOK() {
			return sdk.ErrInternal(err.Error()).Result()
		}
		return res
	}
}

// Handle a message to initialize candidate
func handleMsgInitializeCandidate(ctx sdk.Context, keeper Keeper, msg MsgInitializeCandidate, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return errRes, fmt.Errorf("Failed to query candidate profile: %s", err)
	}

	if !candidateInfo.Initialized {
		return errRes, fmt.Errorf("Candidate has not been initialized")
	}

	accAddress := sdk.AccAddress(candidateInfo.SidechainAddr)
	InitAccount(ctx, keeper, accAddress)

	_, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		log.Infof("Created a new profile for candidate %s account %x", msg.EthAddress, accAddress)
		keeper.SetCandidate(ctx, msg.EthAddress, NewCandidate(accAddress))
	}

	return okRes, nil
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

	pk, err := sdk.GetConsPubKeyBech32(msg.PubKey)
	if err != nil {
		return errRes, fmt.Errorf("GetConsPubKeyBech32 err: %s", err)
	}

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return errRes, fmt.Errorf("Failed to query candidate profile: %s", err)
	}

	if !mainchain.IsBonded(candidateInfo) {
		return errRes, fmt.Errorf("Candidate is not in validator set")
	}

	if !sdk.AccAddress(candidateInfo.SidechainAddr).Equals(msg.Sender) {
		return errRes, fmt.Errorf("Sender has different address recorded on mainchain. mainchain record: %x; sender: %x", sdk.AccAddress(candidateInfo.SidechainAddr), msg.Sender)
	}

	candidate, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		return errRes, fmt.Errorf("Candidate does not exist")
	}

	// Make sure both val address and pub address have not been used before
	valAddress := sdk.ValAddress(candidateInfo.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	_, f := keeper.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk))
	if found != f {
		return errRes, fmt.Errorf("Invalid sender address or public key")
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
	for _, transactor := range candidate.Transactors {
		InitAccount(ctx, keeper, transactor)
	}

	keeper.SetCandidate(ctx, msg.EthAddress, candidate)

	return okRes, nil
}

// Handle a message to sync validator
func handleMsgSyncValidator(ctx sdk.Context, keeper Keeper, msg MsgSyncValidator, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	candidateInfo, err := GetCandidateInfoFromMainchain(ctx, keeper, msg.EthAddress)
	if err != nil {
		return errRes, fmt.Errorf("Failed to query candidate profile: %s", err)
	}

	valAddress := sdk.ValAddress(candidateInfo.SidechainAddr)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return okRes, fmt.Errorf("Validator does not exist")
	}

	updateValidatorToken(ctx, keeper, validator, candidateInfo.StakingPool)
	if !mainchain.IsBonded(candidateInfo) {
		validator.Status = mainchain.ParseStatus(candidateInfo)
		keeper.stakingKeeper.SetValidator(ctx, validator)
		keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	}

	return okRes, nil
}

// Handle a message to sync delegator
func handleMsgSyncDelegator(ctx sdk.Context, keeper Keeper, msg MsgSyncDelegator, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.CandidateAddr = msg.CandidateAddress
	logEntry.DelegatorAddr = msg.DelegatorAddress

	delegator := keeper.GetDelegator(ctx, msg.CandidateAddress, msg.DelegatorAddress)
	di, err := keeper.ethClient.Guard.GetDelegatorInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, mainchain.Hex2Addr(msg.CandidateAddress), mainchain.Hex2Addr(msg.DelegatorAddress))
	if err != nil {
		return errRes, fmt.Errorf("Failed to query delegator info: %s", err)
	}

	delegator.DelegatedStake = sdk.NewIntFromBigInt(di.DelegatedStake)
	keeper.SetDelegator(ctx, msg.CandidateAddress, msg.DelegatorAddress, delegator)
	keeper.SnapshotCandidate(ctx, msg.CandidateAddress)

	return okRes, nil
}

// Handle a message to withdraw reward
func handleMsgWithdrawReward(ctx sdk.Context, keeper Keeper, msg MsgWithdrawReward, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return errRes, fmt.Errorf("Reward does not exist")
	}
	if !reward.HasNewReward() {
		return errRes, fmt.Errorf("No new reward")
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
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// Handle a message to sign reward
func handleMsgSignReward(ctx sdk.Context, keeper Keeper, msg MsgSignReward, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	validator, found := keeper.stakingKeeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return errRes, fmt.Errorf("Sender is not validator")
	}
	if validator.Status != sdk.Bonded {
		return errRes, fmt.Errorf("Validator is not bonded")
	}

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return errRes, fmt.Errorf("Reward does not exist")
	}

	err := reward.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return errRes, fmt.Errorf("Failed to add sig: %s", err)
	}

	keeper.SetReward(ctx, reward)
	return okRes, nil
}

func updateValidatorToken(ctx sdk.Context, keeper Keeper, validator staking.Validator, totalTokens *big.Int) {
	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = sdk.NewIntFromBigInt(totalTokens).Quo(PowerReduction)
	validator.DelegatorShares = validator.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)
	keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
}
