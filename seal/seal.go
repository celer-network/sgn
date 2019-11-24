package seal

import (
	"time"

	"github.com/celer-network/goutils/log"
)

func NewTransactorLog() *TransactorLog {
	msgtypes := make(map[string]uint32)
	now := time.Now().UnixNano()
	return &TransactorLog{
		MsgType:         msgtypes,
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

func NewServiceMsgLog() *ServiceMsgLog {
	now := time.Now().UnixNano()
	return &ServiceMsgLog{
		ExecutionTimeMs: (float64)(now),
	}
}

func CommitServiceMsgLog(entry *ServiceMsgLog) {
	now := time.Now().UnixNano()
	entry.ExecutionTimeMs = ((float64)(now) - entry.ExecutionTimeMs) / 1000000
	if len(entry.Error) > 0 {
		log.Errorln("ServiceMsgLog:", entry)
	} else {
		log.Infoln("ServiceMsgLog:", entry)
	}
}
