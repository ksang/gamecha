package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/ksang/gamecha/seeker"
	"github.com/ksang/gamecha/store"
)

func TestParseSeekerConfig(t *testing.T) {
	var tests = []struct {
		s string
		d seeker.Config
	}{
		{
			`
            seeker:
                steam:
                    portal:  http://api.steampowered.com/
                    key: 16A02FCADCE5D2C8A90CBD9F8A16E63C
                    retry_interval: 0s
            store:
                badger:
`,
			seeker.Config{
				SteamConfig: seeker.SteamConfig{
					Portal:    "http://api.steampowered.com/",
					Key:       "16A02FCADCE5D2C8A90CBD9F8A16E63C",
					WorkerNum: 10,
				},
			},
		},
		{
			`
            seeker:
                steam:
                    portal:  http://api.steampowered.com/
                    key: 16A02FCADCE5D2C8A90CBD9F8A16E63C
                    worker: 10
                    retry_interval: 30s
            store:
                type: bolt`,
			seeker.Config{
				SteamConfig: seeker.SteamConfig{
					Portal:        "http://api.steampowered.com/",
					Key:           "16A02FCADCE5D2C8A90CBD9F8A16E63C",
					WorkerNum:     10,
					RetryInterval: time.Duration(30000000000),
				},
			},
		},
	}

	for caseid, c := range tests {
		res, err := ParseSeekerConfig(c.s)
		if err != nil {
			t.Errorf("case #%d, err: %v", caseid+1, err)
		}
		if !reflect.DeepEqual(res, &c.d) {
			t.Errorf("case #%d, got: %v, expected: %v", caseid+1, res, &c.d)
		}
		t.Logf("Result: %v", res)
	}
}

func TestParseStoreConfig(t *testing.T) {
	var tests = []struct {
		s string
		d store.Config
	}{
		{
			`
            seeker:
                steam:
                    portal:  http://api.steampowered.com/
                    key: 16A02FCADCE5D2C8A90CBD9F8A16E63C
                    retry_interval: 0s
            store:
                type: bolt
                path: /tmp/test.db
`,
			store.Config{
				Database:  "bolt",
				StorePath: "/tmp/test.db",
			},
		},
	}

	for caseid, c := range tests {
		res, err := ParseStoreConfig(c.s)
		if err != nil {
			t.Errorf("case #%d, err: %v", caseid+1, err)
		}
		if !reflect.DeepEqual(res, &c.d) {
			t.Errorf("case #%d, got: %v, expected: %v", caseid+1, res, &c.d)
		}
		t.Logf("Result: %v", res)
	}
}
