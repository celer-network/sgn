// Copyright 2018-2019 Celer Network
//
// This is the Data Access Layer. It maps the server's data structures
// that need to be persisted to KVStore calls:
// * Construct table keys from object attribute(s).
// * Use the appropriate Go data types when fetching stored values.

package watcher

import (
	"errors"
)

const (
	logEventWatch = "lew" // event name -> LogEventID
)

var (
	ErrNoRows = errors.New("No rows matched in the database")
)

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
}

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

func NewDAL(store KVStore) *DAL {
	dal := &DAL{
		st: store,
	}
	return dal
}

func getLogEventWatch(st Storage, name string) (*LogEventID, error) {
	var id LogEventID
	err := st.Get(logEventWatch, name, &id)
	return &id, err
}

func putLogEventWatch(st Storage, name string, id *LogEventID) error {
	return st.Put(logEventWatch, name, id)
}

func deleteLogEventWatch(st Storage, name string) error {
	return st.Delete(logEventWatch, name)
}

func hasLogEventWatch(st Storage, name string) (bool, error) {
	return st.Has(logEventWatch, name)
}

func getAllLogEventWatchKeys(st Storage) ([]string, error) {
	return st.GetKeysByPrefix(logEventWatch, "")
}

func (d *DAL) GetLogEventWatch(name string) (*LogEventID, error) {
	return getLogEventWatch(d.st, name)
}

func (d *DAL) PutLogEventWatch(name string, id *LogEventID) error {
	return putLogEventWatch(d.st, name, id)
}

func (d *DAL) DeleteLogEventWatch(name string) error {
	return deleteLogEventWatch(d.st, name)
}

func (d *DAL) HasLogEventWatch(name string) (bool, error) {
	return hasLogEventWatch(d.st, name)
}

func (d *DAL) GetAllLogEventWatchKeys() ([]string, error) {
	return getAllLogEventWatchKeys(d.st)
}
