package store

import (
	"errors"
	"log"

	"github.com/etcd-io/bbolt"
)

// BoltStore represents a bolt store for game database
type BoltStore struct {
	LogLevel   string
	BucketName string
	db         *bbolt.DB
}

// Close badger store
func (bs *BoltStore) Close() error {
	return nil
}

// SaveGameList to badger store
func (bs *BoltStore) SaveGameList(platform string, games map[int]string) error {
	log.Printf("Saving list %d %s games", len(games), platform)
	value, err := Encode(games)
	if err != nil {
		return err
	}
	key := []byte(platform + "/" + StoreGameListKey)
	if err := bs.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(StoreBucketName))
		if b == nil {
			return errors.New("bolt store no bucket found:" + string(StoreBucketName))
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

// GetGameList from bolt store
func (bs *BoltStore) GetGameList(platform string) (map[int]string, error) {
	var value []byte
	key := []byte(platform + "/" + StoreGameListKey)
	if err := bs.db.View(func(tx *bbolt.Tx) error {
		value = tx.Bucket([]byte(StoreBucketName)).Get(key)
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

// SaveGameRecord to badger store
func (bs *BoltStore) SaveGameRecord(platform string, subid string, r GameRecord) error {
	log.Printf("Saving GameRecord: %s, %s.", platform, r.Name)
	value, err := Encode(r)
	if err != nil {
		return err
	}
	key := []byte(platform + "/" + subid)
	if err := bs.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(StoreBucketName))
		if b == nil {
			return errors.New("bolt store no bucket found:" + string(StoreBucketName))
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
	key := []byte(platform + "/" + subid)
	if err := bs.db.View(func(tx *bbolt.Tx) error {
		value = tx.Bucket([]byte(StoreBucketName)).Get(key)
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

// NewBoltStore creates a bolt store
func NewBoltStore(cfg Config) (*BoltStore, error) {
	db, err := bbolt.Open(cfg.StorePath, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err := db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(StoreBucketName)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	log.Printf("%s store created at: %s bucket: %s", cfg.Database, cfg.StorePath, StoreBucketName)
	return &BoltStore{
		LogLevel:   "debug",
		BucketName: string(StoreBucketName),
		db:         db,
	}, nil
}
