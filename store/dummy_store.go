package store

import (
	"fmt"
	"log"
	"sort"
)

// DummyStore represents a dummy store for game database
type DummyStore struct {
	LogLevel string
}

// Close dummy store
func (ds *DummyStore) Close() error {
	return nil
}

// SaveGameList to dummy store
func (ds *DummyStore) SaveGameList(platform string, games map[int]string) error {
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

// GetGameList from dummy store, pretend to have most of games
func (ds *DummyStore) GetGameList(platform string) (map[int]string, error) {
	ret := make(map[int]string)
	for i := 11; i < 1100000; i++ {
		ret[i] = ""
	}
	return ret, nil
}

// GetSavedGameList from dummy store, pretend to have most of games
func (ds *DummyStore) GetSavedGameList(platform string) (map[int]string, error) {
	ret := make(map[int]string)
	for i := 11; i < 1100000; i++ {
		ret[i] = ""
	}
	return ret, nil
}

// NewDummyStore creates a dummy store
func NewDummyStore(cfg Config) (*DummyStore, error) {
	return &DummyStore{
		LogLevel: "debug",
	}, nil
}
