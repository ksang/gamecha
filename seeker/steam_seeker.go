package seeker

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/ksang/gamecha/store"
	"github.com/valyala/fastjson"
)

var (
	pathGetAppList = "ISteamApps/GetAppList/v2"
)

// SteamConfig is the configuration struct of steam seeker
type SteamConfig struct {
	Portal    string
	Key       string
	ThreadNum int
	client    *http.Client
	store     store.GameStore
}

// GetAppList response parsing
type appListResponse struct {
	Applist appListData `json:"applist"`
}

type appListData struct {
	Apps []appBriefData `json:"apps"`
}

type appBriefData struct {
	Appid int    `json:"appid"`
	Name  string `json:"name"`
}

func startSteamSeeker(ctx context.Context, cfg *SteamConfig) error {
	return nil
}

func (steam *SteamConfig) getSteamAppList(ctx context.Context) error {
	req, err := http.NewRequest("GET", steam.Portal+pathGetAppList, nil)
	if err != nil {
		return err
	}
	return httpDo(ctx, req, steam.client, steam.processSteamAppList)
}

func (steam *SteamConfig) processSteamAppList(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var p fastjson.Parser
	gameList := map[int]string{}
	root, err := p.ParseBytes(body)
	if err != nil {
		return err
	}
	applist := root.Get("applist")
	apps := applist.GetArray("apps")
	for _, game := range apps {
		gameList[game.GetInt("appid")] = string(game.GetStringBytes("name"))
	}
	steam.store.SaveGameList(gameList)
	return nil
}
