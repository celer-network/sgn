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

func NewTransactor(cliHome, chainID, nodeURI, accName, passphrase string, cdc *codec.Codec) (*Transactor, error) {
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
		WithTrustNode(true).
		WithBroadcastMode("sync")

	return &Transactor{
		TxBuilder:  txBldr,
		CliCtx:     cliCtx,
		Key:        key,
		Passphrase: passphrase,
	}, nil
}

func (t *Transactor) BroadcastTx(msg sdk.Msg) (sdk.TxResponse, error) {
	txBldr, err := utils.PrepareTxBuilder(t.TxBuilder, t.CliCtx)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	txBytes, err := txBldr.BuildAndSign(t.Key.GetName(), t.Passphrase, []sdk.Msg{msg})
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return t.CliCtx.BroadcastTx(txBytes)
}
