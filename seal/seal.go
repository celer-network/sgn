package seal

import "github.com/celer-network/goutils/log"

func NewTransactorLog() *TransactorLog {
	msgtypes := make(map[string]uint32)
	return &TransactorLog{
		MsgType: msgtypes,
	}
}

func AddTransactorMsg(txlog *TransactorLog, msgtype string) {
	txlog.MsgType[msgtype] = txlog.MsgType[msgtype] + 1
}

func CommitTransactorLog(txlog *TransactorLog) {
	if len(txlog.Error) > 0 {
		log.Errorln("TransactorLog:", txlog)
	} else {
		log.Infoln("TransactorLog:", txlog)
	}
}
