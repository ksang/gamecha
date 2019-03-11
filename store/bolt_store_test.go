package store

import (
	"testing"
)

var storeCfg = Config{
	Database:  "bolt",
	StorePath: "test.db",
}

func TestNewBoltStore(t *testing.T) {
	store, err := NewBoltStore(storeCfg)
	if err != nil {
		t.Errorf("TestNewBoltStore err: %v", err)
	}
	t.Logf("TestNewBoltStore store: %#v", store)
	t.Logf("TestNewBoltStore db: %#v", store.db)
	if err := store.db.Close(); err != nil {
		t.Errorf("TestNewBoltStore err: %v", err)
	}
}
