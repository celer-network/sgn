package monitor

import (
	"context"
	"fmt"
	"strconv"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client"
	tm "github.com/tendermint/tendermint/types"
)

var (
	initiateWithdrawRewardEvent = fmt.Sprintf("%s.%s='%s'", validator.ModuleName, sdk.AttributeKeyAction, validator.ActionInitiateWithdraw)
	slashEvent                  = fmt.Sprintf("%s.%s='%s'", slash.EventTypeSlash, sdk.AttributeKeyAction, slash.ActionPenalty)
)

func (m *Monitor) monitorSidechainWithdrawReward() {
	m.monitorTendermintEvent(initiateWithdrawRewardEvent, func(e abci.Event) {
		if !m.isValidator {
			return
		}

		event := sdk.StringifyEvent(e)
		if event.Attributes[0].Value == validator.ActionInitiateWithdraw {
			m.handleInitiateWithdrawReward(event.Attributes[1].Value)
		}
	})
}

func (m *Monitor) monitorSidechainSlash() {
	m.monitorTendermintEvent(slashEvent, func(e abci.Event) {
		if !m.isValidator {
			return
		}

		event := sdk.StringifyEvent(e)

		if event.Attributes[0].Value == slash.ActionPenalty {
			nonce, err := strconv.ParseUint(event.Attributes[1].Value, 10, 64)
			if err != nil {
				log.Errorln("Parse penalty nonce error", err)
				return
			}

			penaltyEvent := NewPenaltyEvent(nonce)
			m.handlePenalty(penaltyEvent)
			err = m.db.Set(GetPenaltyKey(penaltyEvent.Nonce), penaltyEvent.MustMarshal())
			if err != nil {
				log.Errorln("db Set err", err)
			}
		}
	})
}

func (m *Monitor) monitorTendermintEvent(eventTag string, handleEvent func(event abci.Event)) {
	client, err := client.NewHTTP(m.operator.CliCtx.NodeURI, "/websocket")
	if err != nil {
		log.Errorln("Fail to start create http client", err)
		return
	}

	err = client.Start()
	if err != nil {
		log.Errorln("Fail to start ws client", err)
		return
	}
	defer client.Stop()

	txs, err := client.Subscribe(context.Background(), "monitor", eventTag)
	if err != nil {
		log.Errorln("ws client subscribe error", err)
		return
	}

	for e := range txs {
		switch data := e.Data.(type) {
		case tm.EventDataNewBlock:
			for _, event := range data.ResultBeginBlock.Events {
				handleEvent(event)
			}
			for _, event := range data.ResultEndBlock.Events {
				handleEvent(event)
			}
		case tm.EventDataTx:
			for _, event := range data.TxResult.Result.Events {
				handleEvent(event)
			}
		}
	}
}

func (m *Monitor) handleInitiateWithdrawReward(ethAddr string) {
	log.Infoln("New initiate withdraw", ethAddr)

	reward, err := validator.CLIQueryReward(m.operator.CliCtx, validator.StoreKey, ethAddr)
	if err != nil {
		log.Errorln("Query reward err", err)
		return
	}

	sig, err := m.ethClient.SignEthMessage(reward.RewardProtoBytes)
	if err != nil {
		log.Errorln("SignEthMessage err", err)
		return
	}

	msg := validator.NewMsgSignReward(ethAddr, sig, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}

func (m *Monitor) handlePenalty(penaltyEvent PenaltyEvent) {
	penalty, err := slash.CLIQueryPenalty(m.operator.CliCtx, slash.StoreKey, penaltyEvent.Nonce)
	if err != nil {
		log.Errorf("Query penalty %d err %s", penaltyEvent.Nonce, err)
		return
	}
	log.Infof("New penalty to %s, reason %s, nonce %d", penalty.ValidatorAddr, penalty.Reason, penaltyEvent.Nonce)

	sig, err := m.ethClient.SignEthMessage(penalty.PenaltyProtoBytes)
	if err != nil {
		log.Errorln("SignEthMessage err", err)
		return
	}

	msg := slash.NewMsgSignPenalty(penaltyEvent.Nonce, sig, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}
