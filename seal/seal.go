package seal

import "github.com/celer-network/goutils/log"

func NewTransactorLog() *TransactorLog {
	msgtypes := make(map[string]uint32)
	return &TransactorLog{
		MsgType: msgtypes,
	}
}

func CommitTransactorLog(txlog *TransactorLog) {
	if len(txlog.Error) > 0 {
		log.Errorln("TransactorLog:", txlog)
	} else {
		log.Infoln("TransactorLog:", txlog)
	}
}

func NewServiceMsgLog() *ServiceMsgLog {
	return &ServiceMsgLog{}
}

func CommitServiceMsgLog(svlog *ServiceMsgLog) {
	if len(svlog.Error) > 0 {
		log.Errorln("ServiceMsgLog:", svlog)
	} else {
		log.Infoln("ServiceMsgLog:", svlog)
	}
}
