package main

import (
	"reflect"
	"testing"
	"github.com/ksang/gamecha/seeker"
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
            store:
                badger:
`,
			seeker.Config{
				SteamConfig: seeker.SteamConfig{
					Portal:    "http://api.steampowered.com/",
					Key:       "16A02FCADCE5D2C8A90CBD9F8A16E63C",
					ThreadNum: 10,
				},
			},
		},
		{
			`
            seeker:
                steam:
                    portal:  http://api.steampowered.com/
                    key: 16A02FCADCE5D2C8A90CBD9F8A16E63C
                    thread: 5
            store:
                badger:
`,
			seeker.Config{
				SteamConfig: seeker.SteamConfig{
					Portal:    "http://api.steampowered.com/",
					Key:       "16A02FCADCE5D2C8A90CBD9F8A16E63C",
					ThreadNum: 5,
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
