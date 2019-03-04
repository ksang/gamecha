package store

// DummyStore represents a dummy store for game database
type DummyStore struct {
	LogLevel string
}

// Close dummy store
func (ds *DummyStore) Close() error {
	return nil
}

// NewDummyStore creates a dummy store
func NewDummyStore(cfg Config) (*DummyStore, error) {
	return &DummyStore{
		LogLevel: "debug",
	}, nil
}
