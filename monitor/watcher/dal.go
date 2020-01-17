// Copyright 2018-2019 Celer Network
//
// This is the Data Access Layer. It maps the server's data structures
// that need to be persisted to KVStore calls:
// * Construct table keys from object attribute(s).
// * Use the appropriate Go data types when fetching stored values.

package watcher

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/celer-network/goutils/log"
)

const (
	// wyre transfer id to status string, client side only
	wyreXferStatus = "wxs" // wyre transfer id -> status

	// OSP only. to avoid handling concurrent openchan requests
	openChannelTs = "oct" // peer@token -> OpenChannelTs

	transactionalMaxRetry   = 10
	transactionalRetryDelay = 10 * time.Millisecond
)

var (
	ErrNoRows = errors.New("No rows matched in the database")
)

// Transaction is the interface implemented by the local and remote stores.
type Transaction interface {
	Commit() error
	Discard()
	ConvertError(err error) error
	Put(table, key string, value interface{}) error
	Get(table, key string, value interface{}) error
	Delete(table, key string) error
	Has(table, key string) (bool, error)
	GetKeysByPrefix(table, prefix string) ([]string, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// KVStore is the interface implemented by the local store (LevelDB
// wrapper) and by the remote store (gRPC calls to a store server).
type KVStore interface {
	Close()
	OpenTransaction() (Transaction, error)
	Put(table, key string, value interface{}) error
	Get(table, key string, value interface{}) error
	Delete(table, key string) error
	Has(table, key string) (bool, error)
	GetKeysByPrefix(table, prefix string) ([]string, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type DAL struct {
	st KVStore
}

type Storage interface {
	Put(table, key string, value interface{}) error
	Get(table, key string, value interface{}) error
	Delete(table, key string) error
	Has(table, key string) (bool, error)
	GetKeysByPrefix(table, prefix string) ([]string, error)
}

type SqlStorage interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func NewDAL(store KVStore) *DAL {
	dal := &DAL{
		st: store,
	}
	return dal
}

func (d *DAL) InsertMonitor(event string, blockNum uint64, blockIdx int64, restart bool) error {
	return insertMonitor(d.st, event, blockNum, blockIdx, restart)
}

func (d *DAL) GetMonitorBlock(event string) (uint64, int64, bool, error) {
	return getMonitorBlock(d.st, event)
}

func (d *DAL) UpdateMonitorBlock(event string, blockNum uint64, blockIdx int64) error {
	return updateMonitorBlock(d.st, event, blockNum, blockIdx)
}

func (d *DAL) UpsertMonitorBlock(event string, blockNum uint64, blockIdx int64, restart bool) error {
	return upsertMonitorBlock(d.st, event, blockNum, blockIdx, restart)
}

func insertMonitor(
	st SqlStorage,
	event string,
	blockNum uint64,
	blockIdx int64,
	restart bool) error {
	q := `INSERT INTO monitor (event, blocknum, blockidx, restart)
		VALUES ($1, $2, $3, $4)`
	res, err := st.Exec(q, event, blockNum, blockIdx, restart)
	return chkExec(res, err, 1, "insertMonitor")
}

func getMonitorBlock(st SqlStorage, event string) (uint64, int64, bool, error) {
	var blockNum uint64
	var blockIdx int64
	q := `SELECT blocknum, blockidx FROM monitor WHERE event = $1`
	err := st.QueryRow(q, event).Scan(&blockNum, &blockIdx)
	found, err := chkQueryRow(err)
	return blockNum, blockIdx, found, err
}

func updateMonitorBlock(
	st SqlStorage,
	event string,
	blockNum uint64,
	blockIdx int64) error {
	q := `UPDATE monitor SET blocknum = $1, blockidx = $2 WHERE event = $3`
	res, err := st.Exec(q, blockNum, blockIdx, event)
	return chkExec(res, err, 1, "updateMonitorBlock")
}

func upsertMonitorBlock(
	st SqlStorage,
	event string,
	blockNum uint64,
	blockIdx int64,
	restart bool) error {
	q := `INSERT INTO monitor (event, blocknum, blockidx, restart)
		VALUES ($1, $2, $3, $4) ON CONFLICT (event) DO UPDATE
		SET blocknum = excluded.blocknum, blockidx = excluded.blockidx`
	res, err := st.Exec(q, event, blockNum, blockIdx, restart)
	return chkExec(res, err, 1, "upsertMonitorBlock")
}

func chkExec(res sql.Result, err error, want int64, caller string) error {
	var got int64
	if err == nil {
		got, err = res.RowsAffected()
		if err == nil && got != want {
			if got == 0 {
				// Wrap ErrNoRows with additional info.
				err = fmt.Errorf("%s: invalid SQL #rows: %d != %d: %w", caller, got, want, ErrNoRows)
			} else {
				err = fmt.Errorf("%s: invalid SQL #rows: %d != %d", caller, got, want)
			}
		}
	}
	return err
}

func chkQueryRow(err error) (bool, error) {
	found := false
	if err == nil {
		found = true
	} else if err == sql.ErrNoRows {
		err = nil
	} else {
		log.Debugln("chkQueryRow SQL error:", err)
	}
	return found, err
}
