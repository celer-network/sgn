package monitor

import (
	"math/big"
	"strings"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/viper"
)

type Operator struct {
	EthClient  *mainchain.EthClient
	Transactor *transactor.Transactor
}

func NewOperator(cdc *codec.Codec, cliHome string) (operator *Operator, err error) {
	ethClient, err := common.NewEthClientFromConfig()
	if err != nil {
		return
	}

	txr, err := transactor.NewTransactor(
		cliHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetString(common.FlagSgnValidatorAccount),
		viper.GetString(common.FlagSgnPassphrase),
		cdc,
		transactor.NewGasPriceEstimator(viper.GetString(common.FlagSgnNodeURI)),
	)
	if err != nil {
		return
	}
	txr.Run()

	return &Operator{
		EthClient:  ethClient,
		Transactor: txr,
	}, nil
}

func (o *Operator) SyncUpdateSidechainAddr(candidateAddr mainchain.Addr) {
	sidechainAddr, err := o.EthClient.SGN.SidechainAddrMap(&bind.CallOpts{}, candidateAddr)
	if err != nil {
		log.Errorln("Query sidechain address error:", err)
		return
	}

	c, err := validator.CLIQueryCandidate(o.Transactor.CliCtx, validator.RouterKey, mainchain.Addr2Hex(candidateAddr))
	if err == nil && sdk.AccAddress(sidechainAddr).Equals(c.ValAccount) {
		log.Debugf("sidechain address of candidate %x is already updated", candidateAddr)
		return
	}

	candidate := validator.NewCandidate(candidateAddr.Hex(), sdk.AccAddress(sidechainAddr))
	candidateData := o.Transactor.CliCtx.Codec.MustMarshalBinaryBare(candidate)
	msg := sync.NewMsgSubmitChange(sync.UpdateSidechainAddr, candidateData, o.EthClient.Client, o.Transactor.Key.GetAddress())
	log.Infof("submit change tx: update sidechain addr for candidate %s %s", candidate.EthAddress, candidate.ValAccount.String())
	o.Transactor.AddTxMsg(msg)
}

// return true if already updated or no need for retry
func (o *Operator) SyncValidator(candidateAddr mainchain.Addr) bool {
	candidate, err := validator.CLIQueryCandidate(o.Transactor.CliCtx, validator.RouterKey, candidateAddr.Hex())
	if err != nil {
		log.Errorln("sidechain query candidate err:", err)
		return false
	}

	candidateInfo, err := o.EthClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, candidateAddr)
	if err != nil {
		log.Errorln("Failed to query candidate info:", err)
		return false
	}

	commission, err := common.NewCommission(o.EthClient, candidateInfo.CommissionRate)
	if err != nil {
		log.Errorln("Failed to create new commission:", err)
		return false
	}

	newVal := staking.Validator{
		Description: staking.Description{
			Identity: mainchain.Addr2Hex(candidateAddr),
		},
		Tokens:     sdk.NewIntFromBigInt(candidateInfo.StakingPool), // not QuoRaw(common.TokenDec) yet
		Status:     mainchain.ParseStatus(candidateInfo),
		Commission: commission,
	}

	var candidateOnly bool
	storedVal, err := validator.CLIQueryValidator(o.Transactor.CliCtx, staking.RouterKey, candidate.ValAccount.String())
	if err != nil {
		// validator not found on sidechain
		if !strings.Contains(err.Error(), common.ErrRecordNotFound.Error()) {
			log.Errorf("CLIQueryValidator %x %s, err: %s", candidateAddr, candidate.ValAccount, err)
			return false
		}
		if candidateInfo.Status.Uint64() == mainchain.Unbonding {
			log.Errorf("Unbonding Candidate %x %s is not found on sidechain ", candidateAddr, candidate.ValAccount)
			return true
		}
		if candidateInfo.Status.Uint64() == mainchain.Unbonded {
			log.Debugf("Candidate %x %s is not a validator on sidechain yet, sync candidate only", candidateAddr, candidate.ValAccount)
			candidateOnly = true
		}
	} else if newVal.Status.Equal(storedVal.Status) &&
		newVal.Tokens.QuoRaw(common.TokenDec).Equal(storedVal.Tokens) &&
		newVal.Tokens.Equal(candidate.StakingPool) &&
		newVal.Commission.Rate.Equal(storedVal.Commission.Rate) {
		log.Debugf("validator %x is already updated", candidateAddr)
		return true
	}

	if candidateOnly {
		if newVal.Tokens.Equal(candidate.StakingPool) && newVal.Commission.Rate.Equal(candidate.CommissionRate) {
			log.Debugf("candidate %x is already updated", candidateAddr)
			return true
		}
	}

	if o.EthClient.Address == candidateAddr {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, viper.GetString(common.FlagSgnPubKey))
		if err != nil {
			log.Errorln("GetConsPubKeyBech32 err:", err)
			return false
		}

		newVal.ConsPubKey = pk
	}

	validatorData := o.Transactor.CliCtx.Codec.MustMarshalBinaryBare(newVal)
	msg := sync.NewMsgSubmitChange(sync.SyncValidator, validatorData, o.EthClient.Client, o.Transactor.Key.GetAddress())
	log.Infof("submit change tx: sync validator %x, tokens %s, status %s, Commission %s",
		candidateAddr, newVal.Tokens, newVal.Status, newVal.Commission.Rate)
	o.Transactor.AddTxMsg(msg)
	return false
}

func (o *Operator) SyncDelegator(candidatorAddr, delegatorAddr mainchain.Addr) {
	di, err := o.EthClient.DPoS.GetDelegatorInfo(&bind.CallOpts{}, candidatorAddr, delegatorAddr)
	if err != nil {
		log.Errorf("Failed to query delegator info: %s", err)
		return
	}

	d, err := validator.CLIQueryDelegator(
		o.Transactor.CliCtx, validator.RouterKey, mainchain.Addr2Hex(candidatorAddr), mainchain.Addr2Hex(delegatorAddr))
	if err == nil {
		if d.DelegatedStake.BigInt().Cmp(di.DelegatedStake) == 0 {
			log.Debugf("delegator %x candidate %x stake %s is already updated", delegatorAddr, candidatorAddr, d.DelegatedStake)
			return
		}
	}

	delegator := validator.NewDelegator(mainchain.Addr2Hex(candidatorAddr), mainchain.Addr2Hex(delegatorAddr))
	delegator.DelegatedStake = sdk.NewIntFromBigInt(di.DelegatedStake)
	delegatorData := o.Transactor.CliCtx.Codec.MustMarshalBinaryBare(delegator)
	msg := sync.NewMsgSubmitChange(sync.SyncDelegator, delegatorData, o.EthClient.Client, o.Transactor.Key.GetAddress())
	log.Infof("submit change tx: sync delegator %x candidate %x stake %s", delegatorAddr, candidatorAddr, delegator.DelegatedStake)
	o.Transactor.AddTxMsg(msg)
}

func (o *Operator) SyncSubscriptionBalance(consumerAddr mainchain.Addr, deposit *big.Int) {
	consumerAddrHex := consumerAddr.Hex()
	depositInt := sdk.NewIntFromBigInt(deposit)
	subscription, err := guard.CLIQuerySubscription(o.Transactor.CliCtx, guard.RouterKey, consumerAddrHex)
	if err == nil {
		if subscription.Deposit.Equal(depositInt) {
			log.Infof("Subscription already updated for %s, deposit %s", consumerAddrHex, deposit)
			return
		}
	}

	subscription = guard.NewSubscription(consumerAddrHex)
	subscription.Deposit = depositInt
	subscriptionData := o.Transactor.CliCtx.Codec.MustMarshalBinaryBare(subscription)
	msg := sync.NewMsgSubmitChange(sync.Subscribe, subscriptionData, o.EthClient.Client, o.Transactor.Key.GetAddress())
	log.Infof("Submit change tx: subscribe ethAddress %s, deposit %s", consumerAddrHex, deposit)
	o.Transactor.AddTxMsg(msg)
}
