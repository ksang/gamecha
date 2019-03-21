// Package store provides database functionalities
package store

// GameStore represents general interfaces of store package.
// Implementations are corresponding to different databases.
type GameStore interface {
	Close() error
	SaveGameList(platform string, games map[int]string) error
	GetGameList(platform string) (map[int]string, error)
	GetSavedGameList(platform string) (map[int]string, error)
}

// Config is the configuration struct of seeker
type Config struct {
	Database  string
	StorePath string
	Buckets   []string
}

var (
	// StoreGameListKey is sub-key name placing full game list of a platform
	StoreGameListKey = "index"
)

// New creates a new GameStore according to configuration
func New(cfg *Config) (GameStore, error) {
	if cfg.Database == "dummy" {
		return NewDummyStore(*cfg)
	}
	if cfg.Database == "bolt" {
		return NewBoltStore(*cfg)
	}
	return nil, nil
}

// GameRecord represents detailed game information
type GameRecord struct {
	Name        string
	RequiredAge int
	Description string
	About       string
	Languages   string
	Developers  []string
	Publishers  []string
}
