package seeker

import (
	"context"
	"testing"
	"time"
)

func TestGetSteamAppList(t *testing.T) {
	steam := &SteamConfig{
		Portal:    "http://api.steampowered.com/",
		Key:       "",
		ThreadNum: 10,
	}
	timeout, _ := time.ParseDuration("10s")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := steam.getSteamAppList(ctx); err != nil {
		t.Errorf("getSteamAppList err: %v", err)
	}
}
