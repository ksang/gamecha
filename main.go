package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	defer func() {
		signal.Stop(sigs)
		cancel()
	}()
	go func() {
		select {
		case <-sigs:
			fmt.Println("Signaled to terminate.")
			cancel()
		case <-ctx.Done():
		}
	}()
	log.Fatal(seeker.Start(ctx, seekerCfg, db))
}
