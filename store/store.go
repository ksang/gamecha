// Package store provides database functionalities
package store

// GameStore represents general interfaces of store package.
// Implementations are corresponding to different databases.
type GameStore interface {
	Close() error
}

// Config is the configuration struct of seeker
type Config struct {
	Database string
}

// New creates a new GameStore according to configuration
func New(cfg Config) (GameStore, error) {
	if cfg.Database == "dummy" {
		return NewDummyStore(cfg)
	}
	return nil, nil
}