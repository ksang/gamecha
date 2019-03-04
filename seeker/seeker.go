// Package seeker is responsible of collecting game data from
// various platforms such as steam and gog
package seeker

import "context"

// Config is the configuration struct of seeker
type Config struct {
	SteamConfig
}

// Start seeker threads according to config
func Start(ctx context.Context, cfg *Config) error {
	if err := startSteamSeeker(ctx, &cfg.SteamConfig); err != nil {
		return err
	}
	return nil
}
