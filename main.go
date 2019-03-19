package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"

	"github.com/ksang/gamecha/seeker"
	"github.com/ksang/gamecha/store"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "d", "gamecha.yml", "config file path.")
}

func main() {
	flag.Parse()
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	seekerCfg, err := ParseSeekerConfig(string(config))
	if err != nil {
		log.Fatal(err)
	}
	storeCfg, err := ParseStoreConfig(string(config))
	if err != nil {
		log.Fatal(err)
	}
	db, err := store.New(storeCfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	log.Fatal(seeker.Start(ctx, seekerCfg, db))
}
