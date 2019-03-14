package store

import (
	"reflect"
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

func TestSaveGetGameList(t *testing.T) {
	store, err := NewBoltStore(storeCfg)
	if err != nil {
		t.Errorf("TestNewBoltStore err: %v", err)
		return
	}
	defer store.db.Close()
	var tests = []struct {
		s map[int]string
	}{
		{
			map[int]string{
				123: "CS",
				122: "Witcher 3",
			},
		},
	}

	for caseid, c := range tests {
		if err := store.SaveGameList("test", c.s); err != nil {
			t.Errorf("case #%d, SaveGameList err: %v", caseid+1, err)
		}
		gl, err := store.GetGameList("test")
		if err != nil {
			t.Errorf("case #%d, decode err: %v", caseid+1, err)
		}
		if !reflect.DeepEqual(gl, c.s) {
			t.Errorf("case #%d, got: %v, expected: %v", caseid+1, gl, c.s)
		}
		t.Logf("Result: %v", gl)
	}
}
