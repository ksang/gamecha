package seeker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	var data appListResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	fmt.Printf("%d apps collected\n", len(data.Applist.Apps))
	return nil
}
