package seal

import (
	"github.com/celer-network/goutils/log"
)

func NewTransactorLog() *TransactorLog {
	msgtypes := make(map[string]uint32)
	return &TransactorLog{
		MsgType: msgtypes,
	}
}

func CommitTransactorLog(entry *TransactorLog) {
	if len(entry.Error) > 0 {
		log.Errorln("TransactorLog:", entry)
	} else {
		log.Infoln("TransactorLog:", entry)
	}
}

func NewServiceMsgLog() *ServiceMsgLog {
	return &ServiceMsgLog{}
}

func CommitServiceMsgLog(entry *ServiceMsgLog) {
	if len(entry.Error) > 0 {
		log.Errorln("ServiceMsgLog:", entry)
	} else {
		log.Infoln("ServiceMsgLog:", entry)
	}
}
