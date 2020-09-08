package monitor

import (
	"encoding/json"
	"fmt"

	tmdb "github.com/tendermint/tm-db"
)

const (
	prefixMonitorWatcher = "emw" // event name -> monitorBlockInfo
)

type watcherDAL struct {
	db tmdb.DB
}

type monitorBlockInfo struct {
	BlockNum uint64 `json:"block_num"`
	BlockIdx int64  `json:"block_idx"`
	Restart  bool   `json:"restart"`
}

func newWatcherDAL(db tmdb.DB) *watcherDAL {
	return &watcherDAL{
		db: tmdb.NewPrefixDB(db, []byte(prefixMonitorWatcher)),
	}
}

func newMonitorBlockInfo(blockNum uint64, blockIdx int64, restart bool) *monitorBlockInfo {
	return &monitorBlockInfo{
		BlockNum: blockNum,
		BlockIdx: blockIdx,
		Restart:  restart,
	}
}

func (dal *watcherDAL) InsertMonitor(event string, blockNum uint64, blockIdx int64, restart bool) error {
	return dal.insertMonitorBlock(event, blockNum, blockIdx, restart)
}

func (dal *watcherDAL) GetMonitorBlock(event string) (uint64, int64, bool, error) {
	info, err := dal.getMonitorBlock(event)
	if err != nil {
		return 0, 0, false, err
	}
	if info == nil {
		return 0, 0, false, nil
	}
	return info.BlockNum, info.BlockIdx, true, nil
}

func (dal *watcherDAL) UpdateMonitorBlock(event string, blockNum uint64, blockIdx int64) error {
	info, err := dal.getMonitorBlock(event)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("MonitorBlockInfo not found")
	}
	return dal.insertMonitorBlock(event, blockNum, blockIdx, info.Restart)
}

func (dal *watcherDAL) UpsertMonitorBlock(event string, blockNum uint64, blockIdx int64, restart bool) error {
	return dal.insertMonitorBlock(event, blockNum, blockIdx, restart)
}

func (dal *watcherDAL) insertMonitorBlock(event string, blockNum uint64, blockIdx int64, restart bool) error {
	info := newMonitorBlockInfo(blockNum, blockIdx, restart)
	bytes, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return dal.db.Set([]byte(event), bytes)
}

func (dal *watcherDAL) getMonitorBlock(event string) (*monitorBlockInfo, error) {
	bytes, err := dal.db.Get([]byte(event))
	if err != nil {
		return nil, err
	}
	if bytes == nil {
		return nil, nil
	}
	var info *monitorBlockInfo = new(monitorBlockInfo)
	err = json.Unmarshal(bytes, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}
