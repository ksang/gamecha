package seeker

import (
	"context"
	"encoding/json"
	"errors"
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

// appdetail response parsing
type steamAppDetailResp map[string]steamAppDetail

type steamAppDetail struct {
	Success bool               `json:"success"`
	Data    steamAppDetailData `json:"data"`
}

type steamAppDetailData struct {
	Typ                 string                 `json:"type"`
	Name                string                 `json:"name"`
	Appid               int                    `json:"steam_appid"`
	RequiredAge         int                    `json:"required_age"`
	IsFree              bool                   `json:"is_free"`
	DetailedDescription string                 `json:"detailed_description"`
	AboutTheGame        string                 `json:"about_the_game"`
	ShortDescription    string                 `json:"short_description"`
	SupportLanguages    string                 `json:"supported_languages"`
	HeaderImage         string                 `json:"header_image"`
	Website             string                 `json:"website"`
	PcRequirements      map[string]string      `json:"pc_requirements"`
	MacRequirements     map[string]string      `json:"mac_requirements"`
	LinuxRequirements   map[string]string      `json:"linux_requirements"`
	Developers          []string               `json:"developers"`
	Publishers          []string               `json:"publishers"`
	PriceOverview       map[string]interface{} `json:"price_overview"`
	Packages            []int                  `json:"packages"`
	// TODO: parse package group details
	PackageGroups      []map[string]interface{} `json:"package_groups"`
	Platforms          map[string]bool          `json:"platforms"`
	MetaCritic         map[string]interface{}   `json:"metacritic"`
	Categories         []map[string]interface{} `json:"categories"`
	Genres             []map[string]interface{} `json:"genres"`
	Screenshots        []map[string]interface{} `json:"screenshots"`
	Recommendations    map[string]interface{}   `json:"recommendations"`
	Achievements       map[string]interface{}   `json:"achievements"`
	ReleaseDate        map[string]interface{}   `json:"release_date"`
	SupportInfo        map[string]interface{}   `json:"support_info"`
	Background         string                   `json:"background"`
	ContentDescriptors map[string]interface{}   `json:"content_descriptors"`
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
	apps := root.Get("applist").GetArray("apps")
	for _, game := range apps {
		gameList[game.GetInt("appid")] = string(game.GetStringBytes("name"))
	}
	steam.store.SaveGameList("steam", gameList)
	return nil
}

func (steam *SteamConfig) parseSteamAppDetail(data []byte) (steamAppDetail, error) {
	var ret steamAppDetailResp
	if err := json.Unmarshal(data, &ret); err != nil {
		return steamAppDetail{}, err
	}
	if len(ret) == 1 {

		for _, v := range ret {
			if !v.Success {
				// response was not success
				break
			}
			return v, nil
		}
	}
	return steamAppDetail{}, errors.New("malformed or failed appdetail response")
}
