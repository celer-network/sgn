package seal

import (
	"time"

	"github.com/celer-network/goutils/log"
)

func NewTransactorLog(sender string) *TransactorLog {
	msgtypes := make(map[string]uint32)
	now := time.Now().UnixNano()
	return &TransactorLog{
		MsgType:         msgtypes,
		Sender:          sender,
		ExecutionTimeMs: (float64)(now),
	}
}

func CommitTransactorLog(entry *TransactorLog) {
	now := time.Now().UnixNano()
	entry.ExecutionTimeMs = ((float64)(now) - entry.ExecutionTimeMs) / 1000000
	if len(entry.Error) > 0 {
		log.Errorln("TransactorLog:", entry)
	} else {
		log.Infoln("TransactorLog:", entry)
	}
}

func NewMsgLog() *MsgLog {
	now := time.Now().UnixNano()
	return &MsgLog{
		ExecutionTimeMs: (float64)(now),
	}
}

func CommitMsgLog(entry *MsgLog) {
	now := time.Now().UnixNano()
	entry.ExecutionTimeMs = ((float64)(now) - entry.ExecutionTimeMs) / 1000000
	if len(entry.Error) > 0 {
		log.Errorln("MsgLog:", entry)
	} else {
		log.Infoln("MsgLog:", entry)
	}
}
