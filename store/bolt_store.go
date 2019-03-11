package store

import (
	"fmt"
	"log"
	"sort"

	"github.com/etcd-io/bbolt"
)

var (
	varBucketName = "gamecha"
)

// BoltStore represents a bolt store for game database
type BoltStore struct {
	LogLevel   string
	BucketName string
	db         *bbolt.DB
}

// Close badger store
func (ds *BoltStore) Close() error {
	return nil
}

// SaveGameList to badger store
func (ds *BoltStore) SaveGameList(platform string, games map[int]string) error {
	log.Printf("Saving %d %s games", len(games), platform)
	var keys []int
	for k := range games {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for i, k := range keys {
		if i < 10 || i >= len(games)-10 {
			fmt.Println("Id:", k, "Name:", games[k])
		}
		if i == 10 && len(games) > 20 {
			fmt.Println("...")
		}
	}
	return nil
}

// GetGameList from bolt store
func (ds *BoltStore) GetGameList(platform string) (map[int]string, error) {
	return make(map[int]string), nil
}

// NewBoltStore creates a bolt store
func NewBoltStore(cfg Config) (*BoltStore, error) {
	db, err := bbolt.Open(cfg.StorePath, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err := db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(varBucketName)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &BoltStore{
		LogLevel:   "debug",
		BucketName: varBucketName,
		db:         db,
	}, nil
}
