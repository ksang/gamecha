package main

import (
	"fmt"

	"github.com/ksang/gamecha/seeker"
	"github.com/olebedev/config"
)

var (
	defaultThreadNum = 10
)

// ParseSeekerConfig parse seeker configurations from string to struct
func ParseSeekerConfig(confStr string) (*seeker.Config, error) {
	cfg, err := config.ParseYaml(confStr)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%#v", cfg)
	seekerConf, err := cfg.Get("seeker")
	if err != nil {
		return nil, err
	}
	steamConf, err := seekerConf.Map("steam")
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", steamConf)
	tn, ok := steamConf["thread"]
	if !ok {
		tn = defaultThreadNum
	}

	return &seeker.Config{
		SteamConfig: seeker.SteamConfig{
			Portal:    steamConf["portal"].(string),
			Key:       steamConf["key"].(string),
			ThreadNum: tn.(int),
		},
	}, nil
}
