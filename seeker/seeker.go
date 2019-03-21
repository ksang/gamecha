// Package seeker is responsible of collecting game data from
// various platforms such as steam and gog
package seeker

import (
	"context"

	"github.com/ksang/gamecha/store"
)

// Config is the configuration struct of seeker
type Config struct {
	SteamConfig
}

// Start seeker threads according to config
func Start(ctx context.Context, cfg *Config, db store.GameStore) error {
	ss, err := startSteamSeeker(ctx, cfg.SteamConfig, db)
	if err != nil {
		return err
	}
	return ss.WaitUntilDone(ctx)
}
