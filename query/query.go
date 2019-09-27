// Package query provides store query functions from cmd
package query

import (
	"fmt"

	"github.com/ksang/gamecha/store"
)

type Querier interface {
	GameList() error
}

func New(db store.GameStore, platform string) Querier {
	return &operator{
		db:       db,
		platform: platform,
	}
}

type operator struct {
	db       store.GameStore
	platform string
}

func (o *operator) GameList() error {
	res, err := o.db.GetGameList(o.platform)
	if err != nil {
		return err
	}
	for _, v := range res {
		fmt.Println(v)
	}
	fmt.Printf("Total %d games.\n", len(res))
	return nil
}
