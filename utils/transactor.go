package utils

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type Transactor struct {
	TxBuilder  types.TxBuilder
	CliCtx     context.CLIContext
	Key        keys.Info
	Passphrase string
}

func NewTransactor(cliHome, accName, chainID, nodeURI string, cdc *codec.Codec) (*Transactor, error) {
	kb, err := client.NewKeyBaseFromDir(cliHome)
	if err != nil {
		return nil, err
	}

	key, err := kb.Get(accName)
	if err != nil {
		return nil, err
	}

	txBldr := auth.
		NewTxBuilderFromCLI().
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID).
		WithKeybase(kb)
	cliCtx := client.
		NewCLIContext().
		WithCodec(cdc).
		WithFromAddress(key.GetAddress()).
		WithFromName(key.GetName()).
		WithNodeURI(nodeURI).
		WithBroadcastMode("sync")
	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return nil, err
	}

	return &Transactor{
		TxBuilder: txBldr,
		CliCtx:    cliCtx,
		Key:       key,
	}, nil
}

func (t *Transactor) BroadcastTx(msg sdk.Msg) (sdk.TxResponse, error) {
	txBytes, err := t.TxBuilder.BuildAndSign(t.Key.GetName(), t.Passphrase, []sdk.Msg{msg})
	if err != nil {
		return sdk.TxResponse{}, err
	}

	res, err := t.CliCtx.BroadcastTx(txBytes)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return res, nil
}
