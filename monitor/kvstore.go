// Copyright 2019 Celer Network
//
// Support the KVStore interface using a SQL database server.

package monitor

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var sqlSchemaCmds = [...]string{
	"CREATE TABLE IF NOT EXISTS keyvals ( key TEXT PRIMARY KEY NOT NULL, tbl TEXT NOT NULL, val BYTEA NOT NULL );",
	"CREATE INDEX kvs_tbl_idx ON keyvals (tbl);",
	"CREATE TABLE IF NOT EXISTS channels ( cid TEXT PRIMARY KEY NOT NULL, peer TEXT NOT NULL, token TEXT NOT NULL, ledger TEXT NOT NULL, state INT NOT NULL, statets TIMESTAMPTZ NOT NULL, opents TIMESTAMPTZ NOT NULL, openresp BYTEA, onchainbalance BYTEA, basesn INT NOT NULL, lastusedsn INT NOT NULL, lastackedsn INT NOT NULL, lastnackedsn INT NOT NULL, selfsimplex BYTEA, peersimplex BYTEA, UNIQUE (peer, token) );",
	"CREATE INDEX chan_ledger_idx ON channels (ledger);",
	"CREATE TABLE IF NOT EXISTS closedchannels ( cid TEXT PRIMARY KEY NOT NULL, peer TEXT NOT NULL, token TEXT NOT NULL, opents TIMESTAMPTZ NOT NULL, closets TIMESTAMPTZ NOT NULL );",
	"CREATE INDEX cc_peer_token_idx ON closedchannels (peer, token);",
	"CREATE TABLE IF NOT EXISTS payments ( payid TEXT PRIMARY KEY NOT NULL, pay BYTEA, paynote BYTEA, incid TEXT NOT NULL, instate INT NOT NULL, outcid TEXT NOT NULL, outstate INT NOT NULL, src TEXT NOT NULL, dest TEXT NOT NULL, createts TIMESTAMPTZ NOT NULL );",
	"CREATE INDEX pay_src_idx ON payments (src);",
	"CREATE INDEX pay_dest_idx ON payments (dest);",
	"CREATE INDEX pay_ts_idx ON payments (createts);",
	"CREATE TABLE IF NOT EXISTS paydelegation ( payid TEXT PRIMARY KEY NOT NULL REFERENCES payments (payid) ON UPDATE CASCADE ON DELETE CASCADE, dest TEXT NOT NULL, status INT NOT NULL, payidout TEXT, delegator TEXT );",
	"CREATE INDEX paydel_dest_idx ON paydelegation (dest);",
	"CREATE TABLE IF NOT EXISTS secrets ( hash TEXT PRIMARY KEY NOT NULL, preimage TEXT NOT NULL, payid TEXT NOT NULL, UNIQUE (hash, payid) );",
	"CREATE TABLE IF NOT EXISTS tcb ( addr TEXT NOT NULL, token TEXT NOT NULL, deposit TEXT NOT NULL, UNIQUE (addr, token) );",
	"CREATE TABLE IF NOT EXISTS monitor ( event TEXT PRIMARY KEY NOT NULL, blocknum INT NOT NULL, blockidx INT NOT NULL, restart BOOL NOT NULL );",
	"CREATE TABLE IF NOT EXISTS routing ( dest TEXT NOT NULL, token TEXT NOT NULL, cid TEXT NOT NULL, UNIQUE (dest, token) );",
	"CREATE TABLE IF NOT EXISTS osps ( addr TEXT PRIMARY KEY NOT NULL );",
	"CREATE TABLE IF NOT EXISTS edges ( token TEXT NOT NULL, cid TEXT NOT NULL, addr1 TEXT NOT NULL, addr2 TEXT NOT NULL, UNIQUE (token, cid) );",
	"CREATE TABLE IF NOT EXISTS peers ( peer TEXT PRIMARY KEY NOT NULL, server TEXT NOT NULL, activecids TEXT NOT NULL,  delegateproof BYTEA );",
	"CREATE TABLE IF NOT EXISTS desttokens ( dest TEXT NOT NULL, token TEXT NOT NULL, osps TEXT NOT NULL,  openchanblknum INT NOT NULL, UNIQUE (dest, token) );",
	"CREATE TABLE IF NOT EXISTS chanmessages ( cid TEXT NOT NULL, seqnum INT NOT NULL, msg BYTEA, UNIQUE (cid, seqnum) );",
	"CREATE TABLE IF NOT EXISTS chanmigration ( cid TEXT NOT NULL REFERENCES channels (cid) ON DELETE CASCADE, toledger TEXT NOT NULL, deadline INT NOT NULL, onchainreq BYTEA, state INT NOT NULL, ts TIMESTAMPTZ NOT NULL, UNIQUE (cid, toledger) );",
	"CREATE INDEX mg_toledger_state_idx ON chanmigration (toledger, state);",
}

var (
	ErrNilValue   = errors.New("Value cannot be nil")
	ErrTxConflict = errors.New("Transaction conflict")
)

const (
	separator = "|" // reserved character for keys construction
)

type KVStoreSQL struct {
	driver string  // database driver
	info   string  // database connection info
	crdb   bool    // database is CockroachDB
	db     *sql.DB // database access object
}

type TransactionSQL struct {
	store *KVStoreSQL // remote store handle
	dbTx  *sql.Tx     // database transaction
}

type dbOrTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Create a new remote K/V store.
func NewKVStoreSQL(driver, info string) (*KVStoreSQL, error) {
	s := &KVStoreSQL{
		driver: driver,
		info:   info,
		crdb:   true,
	}

	// Special check for SQLite on the client: if the file
	// does not already exist, then initialize its schema.
	// Note: in the "sqlite3" case "info" is the file path.
	initSchema := false
	if driver == "sqlite3" {
		s.crdb = false
		if ok, err := exists(info); err != nil {
			return nil, err
		} else if !ok {
			initSchema = true
			dir := path.Dir(info)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}

	db, err := sql.Open(driver, info)
	if err != nil {
		return nil, err
	}

	s.db = db

	// Initialize the database schema if needed.
	if initSchema {
		for _, cmd := range sqlSchemaCmds {
			_, err = db.Exec(cmd)
			if err != nil {
				db.Close()
				return nil, err
			}
		}
	}

	return s, nil
}

// Close the remote K/V store.
func (s *KVStoreSQL) Close() {
	if s.db != nil {
		s.db.Close()
		s.db = nil
		s.driver = ""
		s.info = ""
	}
}

func (s *KVStoreSQL) put(db dbOrTx, table, key string, value interface{}) error {
	if err := checkTableKey(table, key); err != nil {
		return err
	}
	if value == nil {
		return ErrNilValue
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	q := `INSERT INTO keyvals (key, tbl, val) VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET val = excluded.val`
	_, err = db.Exec(q, storeKey(table, key), table, data)
	return err
}

func (s *KVStoreSQL) get(db dbOrTx, table, key string, value interface{}) error {
	if err := checkTableKey(table, key); err != nil {
		return err
	}
	if value == nil {
		return ErrNilValue
	}

	var data []byte
	q := "SELECT val FROM keyvals WHERE key = $1"
	err := db.QueryRow(q, storeKey(table, key)).Scan(&data)
	if err == nil {
		err = json.Unmarshal(data, value)
	}
	return err
}

func (s *KVStoreSQL) del(db dbOrTx, table, key string) error {
	if err := checkTableKey(table, key); err != nil {
		return err
	}

	q := "DELETE FROM keyvals WHERE key = $1"
	_, err := db.Exec(q, storeKey(table, key))
	return err
}

func (s *KVStoreSQL) has(db dbOrTx, table, key string) (bool, error) {
	if err := checkTableKey(table, key); err != nil {
		return false, err
	}

	var data int
	q := "SELECT 1 FROM keyvals WHERE key = $1"
	err := db.QueryRow(q, storeKey(table, key)).Scan(&data)
	if err == nil {
		return true, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	}
	return false, err
}

func (s *KVStoreSQL) getKeys(db dbOrTx, table, prefix string) ([]string, error) {
	if err := checkTableKey(table, " "); err != nil {
		return nil, err
	}

	var params []interface{}

	// For an empty prefix the query uses only the indexed "tbl" table.
	q := "SELECT key FROM keyvals WHERE tbl = $1"
	params = append(params, table)

	if prefix != "" {
		// Further filtering on the keys using LIKE prefix matching.
		q += " AND key LIKE $2"
		like := storeKey(table, prefix) + "%"
		params = append(params, like)
	}
	q += " ORDER BY key"

	rows, err := db.Query(q, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err = rows.Scan(&key); err != nil {
			return nil, err
		}

		_, key = tableKey([]byte(key))
		keys = append(keys, key)
	}

	return keys, nil
}

func (s *KVStoreSQL) Put(table, key string, value interface{}) error {
	return s.put(s.db, table, key, value)
}

func (s *KVStoreSQL) Get(table, key string, value interface{}) error {
	return s.get(s.db, table, key, value)
}

func (s *KVStoreSQL) Delete(table, key string) error {
	return s.del(s.db, table, key)
}

func (s *KVStoreSQL) Has(table, key string) (bool, error) {
	return s.has(s.db, table, key)
}

func (s *KVStoreSQL) GetKeysByPrefix(table, prefix string) ([]string, error) {
	return s.getKeys(s.db, table, prefix)
}

func (s *KVStoreSQL) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

func (s *KVStoreSQL) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

func (s *KVStoreSQL) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

func (s *KVStoreSQL) OpenTransaction() (Transaction, error) {
	dbTx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	if s.crdb {
		_, err = dbTx.Exec("SAVEPOINT cockroach_restart")
		if err != nil {
			dbTx.Rollback()
			return nil, err
		}
	}

	tx := &TransactionSQL{
		store: s,
		dbTx:  dbTx,
	}
	return tx, nil
}

func (tx *TransactionSQL) Discard() {
	if tx.dbTx != nil {
		err := tx.dbTx.Rollback()
		if err == nil {
			tx.dbTx = nil
		}
	}
}

func (tx *TransactionSQL) ConvertError(err error) error {
	if err == nil {
		return nil
	}

	// Special re-mapping of this error back to transaction conflict.
	var patterns []string
	if tx.store.crdb {
		patterns = []string{"retry transaction", "restart transaction",
			"current transaction is aborted", "40001", "cr000"}
	} else {
		patterns = []string{"database is locked"}
	}

	errMsg := strings.ToLower(err.Error())
	for _, pat := range patterns {
		if strings.Contains(errMsg, pat) {
			return ErrTxConflict
		}
	}

	return err
}

func (tx *TransactionSQL) Commit() error {
	var err error
	if tx.store.crdb {
		// For CockroachDB, both "release savepoint" and the follow-up
		// "commit" may fail.  The commit after a successful "release"
		// is not a NOP.
		_, err = tx.dbTx.Exec("RELEASE SAVEPOINT cockroach_restart")
	}

	if err == nil {
		err = tx.dbTx.Commit()
		if err == nil {
			tx.dbTx = nil
			return nil
		}
	}

	return tx.ConvertError(err)
}

func (tx *TransactionSQL) Put(table, key string, value interface{}) error {
	return tx.store.put(tx.dbTx, table, key, value)
}

func (tx *TransactionSQL) Get(table, key string, value interface{}) error {
	return tx.store.get(tx.dbTx, table, key, value)
}

func (tx *TransactionSQL) Delete(table, key string) error {
	return tx.store.del(tx.dbTx, table, key)
}

func (tx *TransactionSQL) Has(table, key string) (bool, error) {
	return tx.store.has(tx.dbTx, table, key)
}

func (tx *TransactionSQL) GetKeysByPrefix(table, prefix string) ([]string, error) {
	return tx.store.getKeys(tx.dbTx, table, prefix)
}

func (tx *TransactionSQL) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.dbTx.Exec(query, args...)
}

func (tx *TransactionSQL) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.dbTx.Query(query, args...)
}

func (tx *TransactionSQL) QueryRow(query string, args ...interface{}) *sql.Row {
	return tx.dbTx.QueryRow(query, args...)
}

// Check if the table and key parameters are valid.
func checkTableKey(table, key string) error {
	if table == "" || key == "" {
		return fmt.Errorf("table and key parameters must be specified")
	}

	// The separator character cannot be used in the table name.
	if strings.Contains(table, separator) {
		return fmt.Errorf("invalid table name: %s", table)
	}
	return nil
}

// storeKey returns the store's key for the given table and entry key.
func storeKey(table, key string) string {
	return table + separator + key
}

// tableKey returns the user visible (table, key) info from a store key.
func tableKey(skey []byte) (string, string) {
	parts := strings.SplitN(string(skey), separator, 2)
	return parts[0], parts[1]
}

func exists(fpath string) (bool, error) {
	_, err := os.Stat(fpath)
	if err == nil || os.IsExist(err) {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
