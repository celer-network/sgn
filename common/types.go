package common

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/viper"
)

type ParamChange struct {
	Record   sdk.Int `json:"record"`
	NewValue sdk.Int `json:"new_value"`
}

func NewParamChange(record, newValue sdk.Int) ParamChange {
	return ParamChange{
		Record:   record,
		NewValue: newValue,
	}
}

// implement fmt.Stringer
func (p ParamChange) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Record: %v, NewValue: %v`, p.Record, p.NewValue))
}

type Sig struct {
	Signer string `json:"signer"`
	Sig    []byte `json:"sig"`
}

func NewSig(signer string, sig []byte) Sig {
	return Sig{
		Signer: signer,
		Sig:    sig,
	}
}

// implement fmt.Stringer
func (r Sig) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Signer: %s, Sig: %x,`, r.Signer, r.Sig))
}

func AddSig(sigs []Sig, msg []byte, sig []byte, expectedSigner string) ([]Sig, error) {
	signer, err := eth.RecoverSigner(msg, sig)
	if err != nil {
		return nil, err
	}

	signerAddr := mainchain.Addr2Hex(signer)
	if signerAddr != mainchain.FormatAddrHex(expectedSigner) {
		err = fmt.Errorf("invalid signer address %s %s", signerAddr, expectedSigner)
		return nil, err
	}

	for i, s := range sigs {
		if s.Signer == signerAddr {
			if bytes.Compare(s.Sig, sig) == 0 {
				// already signed with the same sig
				return sigs, nil
			}
			log.Debugf("repeated signer %s overwite existing sig", signerAddr)
			sigs[i] = NewSig(signerAddr, sig)
			return sigs, nil
		}
	}

	return append(sigs, NewSig(signerAddr, sig)), nil
}

func NewCommission(ethClient *mainchain.EthClient, commissionRate *big.Int) (staking.Commission, error) {
	commissionBase, err := ethClient.DPoS.COMMISSIONRATEBASE(&bind.CallOpts{})
	if err != nil {
		return staking.Commission{}, err
	}

	prec := int64(len(commissionBase.String()) - 1)
	return staking.Commission{
		CommissionRates: staking.CommissionRates{
			Rate:          sdk.NewDecFromBigIntWithPrec(commissionRate, prec),
			MaxRate:       sdk.NewDec(1),
			MaxChangeRate: sdk.NewDec(1),
		},
	}, nil
}

func NewEthClientFromConfig() (*mainchain.EthClient, error) {
	return mainchain.NewEthClient(
		viper.GetString(FlagEthGateway),
		viper.GetString(FlagEthKeystore),
		viper.GetString(FlagEthPassphrase),
		&mainchain.TransactorConfig{
			BlockDelay:           viper.GetUint64(FlagEthBlockDelay),
			BlockPollingInterval: viper.GetUint64(FlagEthPollInterval),
			ChainId:              big.NewInt(viper.GetInt64(FlagEthChainID)),
			AddGasPriceGwei:      viper.GetUint64(FlagEthAddGasPriceGwei),
			MinGasPriceGwei:      viper.GetUint64(FlagEthMinGasPriceGwei),
		},
		viper.GetString(FlagEthDPoSAddress),
		viper.GetString(FlagEthSGNAddress),
	)
}
