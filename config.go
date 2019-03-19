package main

import (
	"time"

	"github.com/ksang/gamecha/seeker"
	"github.com/ksang/gamecha/store"
	"github.com/olebedev/config"
)

var (
	defaultWorkerNum     = 10
	defaultRetryInterval = "30s"
	defaultRetryCount    = 5
	defaultStoreType     = "bolt"
	defaultStorePath     = "gamecha.db"
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
	wn, ok := steamConf["worker"]
	if !ok {
		wn = defaultWorkerNum
	}
	ris, ok := steamConf["retry_interval"]
	if !ok {
		ris = defaultRetryInterval
	}
	rc, ok := steamConf["retry_count"]
	if !ok {
		ris = defaultRetryCount
	}
	ri, err := time.ParseDuration(ris.(string))
	if err != nil {
		return nil, err
	}

	return &seeker.Config{
		SteamConfig: seeker.SteamConfig{
			Portal:        steamConf["portal"].(string),
			Key:           steamConf["key"].(string),
			WorkerNum:     wn.(int),
			RetryInterval: ri,
			RetryCount:    rc.(int),
		},
	}, nil
}

// ParseStoreConfig parse store configurations from string to struct
func ParseStoreConfig(confStr string) (*store.Config, error) {
	cfg, err := config.ParseYaml(confStr)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%#v", cfg)
	storeConf, err := cfg.Map("store")
	if err != nil {
		return nil, err
	}
	typ, ok := storeConf["type"]
	if !ok {
		typ = defaultStoreType
	}
	path, ok := storeConf["path"]
	if !ok {
		path = defaultStorePath
	}
	return &store.Config{
		Database:  typ.(string),
		StorePath: path.(string),
	}, nil
}
