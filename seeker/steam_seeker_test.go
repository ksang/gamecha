package seeker

import (
	"context"
	"testing"
	"time"

	"github.com/ksang/gamecha/store"
)

func TestGetSteamAppList(t *testing.T) {
	s, _ := store.NewDummyStore(store.Config{})
	steam := &SteamConfig{
		Portal:    "http://api.steampowered.com/",
		Key:       "",
		ThreadNum: 10,
		store:     s,
	}
	timeout, _ := time.ParseDuration("60s")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := steam.getSteamAppList(ctx); err != nil {
		t.Errorf("getSteamAppList err: %v", err)
	}

}
