package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ksang/gamecha/query"
	"github.com/ksang/gamecha/seeker"
	"github.com/ksang/gamecha/store"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app        = kingpin.New("gamecha", "Game metadata toolkits.")
	cf         = app.Flag("config", "config file path").Default("gamecha.yml").String()
	sk         = app.Command("seeker", "Start gamecha in seeker mode.")
	op         = app.Command("query", "Query gamecha store.")
	opList     = op.Command("list", "List all games in store.")
	opPlatform = opList.Flag("platform", "Which platform to query").Default("steam").String()
)

func openStore(confStr string) store.GameStore {
	storeCfg, err := ParseStoreConfig(confStr)
	if err != nil {
		log.Fatal(err)
	}
	db, err := store.New(storeCfg)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func startSeeker(cfg string) {
	config, err := ioutil.ReadFile(cfg)
	if err != nil {
		log.Fatal(err)
	}
	db := openStore(string(config))
	seekerCfg, err := ParseSeekerConfig(string(config))
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

func newQuery(cfg string, platform string) query.Querier {
	config, err := ioutil.ReadFile(cfg)
	if err != nil {
		log.Fatal(err)
	}
	db := openStore(string(config))
	return query.New(db, platform)
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Register user
	case "seeker":
		startSeeker(*cf)

		// Post message
	case opList.FullCommand():
		if err := newQuery(*cf, *opPlatform).GameList(); err != nil {
			log.Fatal(err)
		}
	}
}
