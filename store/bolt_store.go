package store

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/etcd-io/bbolt"
)

// BoltStore represents a bolt store for game database
type BoltStore struct {
	LogLevel string
	Buckets  []string
	db       *bbolt.DB
	debugLog *log.Logger
	infoLog  *log.Logger
}

// NewBoltStore creates a bolt store
func NewBoltStore(cfg Config) (*BoltStore, error) {
	db, err := bbolt.Open(cfg.StorePath, 0600, nil)
	if err != nil {
		return nil, err
	}
	for _, b := range cfg.Buckets {
		if err := db.Update(func(tx *bbolt.Tx) error {
			if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	log.Printf("%s store created at: %s buckets: %s", cfg.Database, cfg.StorePath, cfg.Buckets)
	return &BoltStore{
		LogLevel: "debug",
		Buckets:  cfg.Buckets,
		db:       db,
		debugLog: log.New(os.Stdout, "BoltStore DEBUG:", log.LstdFlags|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "BoltStore INFO:", log.LstdFlags|log.Lshortfile),
	}, nil
}

// Close badger store
func (bs *BoltStore) Close() error {
	return nil
}

// SaveGameList to badger store
func (bs *BoltStore) SaveGameList(platform string, games map[int]string) error {
	bs.infoLog.Printf("Saving list %d %s games", len(games), platform)
	value, err := Encode(games)
	if err != nil {
		return err
	}
	key := []byte(StoreGameListKey)
	if err := bs.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(platform))
		if b == nil {
			return errors.New("bolt store no bucket found:" + string(platform))
		}
		if err := b.Put(key, value); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetGameList index from bolt store
func (bs *BoltStore) GetGameList(platform string) (map[int]string, error) {
	var value []byte
	key := []byte(StoreGameListKey)
	if err := bs.db.View(func(tx *bbolt.Tx) error {
		value = tx.Bucket([]byte(platform)).Get(key)
		return nil
	}); err != nil {
		return nil, err
	}
	var games map[int]string
	if len(value) > 0 {
		if err := Decode(value, &games); err != nil {
			return nil, err
		}
	}
	return games, nil
}

// GetSavedGameList from bolt store
func (bs *BoltStore) GetSavedGameList(platform string) (map[int]string, error) {
	games := make(map[int]string)
	if err := bs.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket([]byte(platform)).Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if id, err := strconv.ParseInt(string(k), 10, 32); err == nil {
				games[int(id)] = ""
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return games, nil
}

// SaveGameRecord to badger store
func (bs *BoltStore) SaveGameRecord(platform string, subid string, r GameRecord) error {
	bs.debugLog.Printf("Saving GameRecord: %s, %s - %s.", platform, subid, r.Name)
	value, err := Encode(r)
	if err != nil {
		return err
	}
	key := []byte(subid)
	if err := bs.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(platform))
		if b == nil {
			return errors.New("bolt store no bucket found:" + string(platform))
		}
		if err := b.Put(key, value); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetGameRecord from bolt store
func (bs *BoltStore) GetGameRecord(platform string, subid string) (*GameRecord, error) {
	var value []byte
	key := []byte(subid)
	if err := bs.db.View(func(tx *bbolt.Tx) error {
		value = tx.Bucket([]byte(platform)).Get(key)
		return nil
	}); err != nil {
		return nil, err
	}
	var r GameRecord
	if len(value) > 0 {
		if err := Decode(value, &r); err != nil {
			return nil, err
		}
	}
	return &r, nil
}
