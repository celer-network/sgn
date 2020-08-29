package sync

import (
	"context"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewMsgSubmitChange(
	changeType string, data []byte, ethClient *ethclient.Client, sender sdk.AccAddress) types.MsgSubmitChange {

	var blkNum uint64
	head, err := ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Errorln("cannot fetch mainchain block number:", err)
	} else {
		blkNum = head.Number.Uint64()
	}
	return types.NewMsgSubmitChange(changeType, data, blkNum, sender)
}
