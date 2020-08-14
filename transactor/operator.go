package transactor

import (
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
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
	Transactor *Transactor
}

func NewOperator(cdc *codec.Codec) (operator *Operator, err error) {
	ethClient, err := mainchain.NewEthClient(
		viper.GetString(common.FlagEthGateway),
		viper.GetString(common.FlagEthKeystore),
		viper.GetString(common.FlagEthPassphrase),
		&mainchain.TransactorConfig{
			BlockDelay:           viper.GetUint64(common.FlagEthBlockDelay),
			BlockPollingInterval: viper.GetUint64(common.FlagEthPollInterval),
			ChainId:              big.NewInt(viper.GetInt64(common.FlagEthChainID)),
			AddGasPriceGwei:      viper.GetUint64(common.FlagEthAddGasPriceGwei),
			MinGasPriceGwei:      viper.GetUint64(common.FlagEthMinGasPriceGwei),
		},
		viper.GetString(common.FlagEthDPoSAddress),
		viper.GetString(common.FlagEthSGNAddress),
		viper.GetString(common.FlagEthLedgerAddress),
	)
	if err != nil {
		return
	}

	transactor, err := NewTransactor(
		viper.GetString(common.FlagCLIHome),
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetString(common.FlagSgnOperator),
		viper.GetString(common.FlagSgnPassphrase),
		cdc,
		NewGasPriceEstimator(viper.GetString(common.FlagSgnNodeURI)),
	)
	if err != nil {
		return
	}

	return &Operator{
		EthClient:  ethClient,
		Transactor: transactor,
	}, nil
}

func (o *Operator) SyncUpdateSidechainAddr(candidateAddr mainchain.Addr) {
	sidechainAddr, err := o.EthClient.SGN.SidechainAddrMap(&bind.CallOpts{}, candidateAddr)
	if err != nil {
		log.Errorln("Query sidechain address error:", err)
		return
	}

	c, err := validator.CLIQueryCandidate(o.Transactor.CliCtx, validator.RouterKey, mainchain.Addr2Hex(candidateAddr))
	if err == nil && sdk.AccAddress(sidechainAddr).Equals(c.Operator) {
		log.Infof("The sidechain address of candidate %x has been updated", candidateAddr)
		return
	}

	candidate := validator.NewCandidate(candidateAddr.Hex(), sdk.AccAddress(sidechainAddr))
	candidateData := o.Transactor.CliCtx.Codec.MustMarshalBinaryBare(candidate)
	msg := sync.NewMsgSubmitChange(sync.UpdateSidechainAddr, candidateData, o.Transactor.Key.GetAddress())
	log.Infof("submit change tx: update sidechain addr for candidate %s %s", candidate.EthAddress, candidate.Operator.String())
	o.Transactor.AddTxMsg(msg)
}

func (o *Operator) SyncValidator(candidateAddr mainchain.Addr) {
	ci, err := o.EthClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, candidateAddr)
	if err != nil {
		log.Errorln("Failed to query candidate info:", err)
		return
	}

	commission, err := common.NewCommission(o.EthClient, ci.CommissionRate)
	if err != nil {
		log.Errorln("Failed to create new commission:", err)
		return
	}

	vt := staking.Validator{
		Description: staking.Description{
			Identity: candidateAddr.Hex(),
		},
		Tokens:     sdk.NewIntFromBigInt(ci.StakingPool).QuoRaw(common.TokenDec),
		Status:     mainchain.ParseStatus(ci),
		Commission: commission,
	}

	candidate, err := validator.CLIQueryCandidate(o.Transactor.CliCtx, validator.RouterKey, candidateAddr.Hex())
	if err != nil {
		log.Errorln("sidechain query candidate err:", err)
		return
	}
	v, err := validator.CLIQueryValidator(
		o.Transactor.CliCtx, staking.RouterKey, candidate.Operator.String())
	if err == nil {
		if vt.Status.Equal(v.Status) && vt.Tokens.Equal(v.Tokens) &&
			vt.Commission.CommissionRates.Rate.Equal(v.Commission.CommissionRates.Rate) {
			log.Infof("no need to sync updated validator %x", candidateAddr)
			return
		}
	}

	if o.EthClient.Address == candidateAddr {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, viper.GetString(common.FlagSgnPubKey))
		if err != nil {
			log.Errorln("GetConsPubKeyBech32 err:", err)
			return
		}

		vt.ConsPubKey = pk
	}

	validatorData := o.Transactor.CliCtx.Codec.MustMarshalBinaryBare(vt)
	msg := sync.NewMsgSubmitChange(sync.SyncValidator, validatorData, o.Transactor.Key.GetAddress())
	log.Infof("submit change tx: sync validator %x", candidateAddr)
	o.Transactor.AddTxMsg(msg)
}

func (o *Operator) SyncDelegator(candidatorAddr, delegatorAddr mainchain.Addr) {
	di, err := o.EthClient.DPoS.GetDelegatorInfo(&bind.CallOpts{}, candidatorAddr, delegatorAddr)
	if err != nil {
		log.Errorf("Failed to query delegator info: %s", err)
		return
	}

	delegator := validator.NewDelegator(mainchain.Addr2Hex(candidatorAddr), mainchain.Addr2Hex(delegatorAddr))
	delegator.DelegatedStake = sdk.NewIntFromBigInt(di.DelegatedStake)
	delegatorData := o.Transactor.CliCtx.Codec.MustMarshalBinaryBare(delegator)
	msg := sync.NewMsgSubmitChange(sync.SyncDelegator, delegatorData, o.Transactor.Key.GetAddress())
	log.Infof("submit change tx: sync delegator %x candidate %x stake %s", delegatorAddr, candidatorAddr, delegator.DelegatedStake)
	o.Transactor.AddTxMsg(msg)
}
