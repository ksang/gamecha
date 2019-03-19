package seeker

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ksang/gamecha/store"
	"github.com/valyala/fastjson"
)

// ContextKey is used to define key used in context
type ContextKey string

var (
	pathGetAppList   = "/ISteamApps/GetAppList/v2"
	pathGetAppDetail = "https://store.steampowered.com/api/appdetails"
	workerIDKey      = ContextKey("workerID")
)

var (
	// ErrSteamFailReponse indicates steam api response with failed message
	ErrSteamFailReponse = errors.New("failed steam response")
	// ErrSteamRateLimit indicates steam api is rate limiting our seeker
	ErrSteamRateLimit = errors.New("steam rate limit")
)

// SteamConfig is the configuration struct of steam seeker
type SteamConfig struct {
	Portal        string
	Key           string
	WorkerNum     int
	RetryInterval time.Duration
	RetryCount    int
}

// SteamSeeker object
type SteamSeeker struct {
	config       SteamConfig
	client       *http.Client
	store        store.GameStore
	queue        chan int
	workerQuit   chan struct{}
	workerDone   chan struct{}
	workerReturn chan store.GameRecord
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
	RequiredAge         interface{}            `json:"required_age"`
	IsFree              bool                   `json:"is_free"`
	DetailedDescription string                 `json:"detailed_description"`
	AboutTheGame        string                 `json:"about_the_game"`
	ShortDescription    string                 `json:"short_description"`
	SupportLanguages    string                 `json:"supported_languages"`
	HeaderImage         string                 `json:"header_image"`
	Website             string                 `json:"website"`
	PcRequirements      interface{}            `json:"pc_requirements"`
	MacRequirements     interface{}            `json:"mac_requirements"`
	LinuxRequirements   interface{}            `json:"linux_requirements"`
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

func startSteamSeeker(ctx context.Context, cfg SteamConfig, db store.GameStore) (*SteamSeeker, error) {
	steam := SteamSeeker{
		config:       cfg,
		queue:        make(chan int),
		store:        db,
		workerReturn: make(chan store.GameRecord, cfg.WorkerNum),
		workerQuit:   make(chan struct{}, cfg.WorkerNum),
		workerDone:   make(chan struct{}, cfg.WorkerNum),
	}
	if err := steam.getSteamAppList(ctx); err != nil {
		return nil, err
	}
	for i := 0; i < cfg.WorkerNum; i++ {
		wCtx := context.WithValue(ctx, workerIDKey, i)
		go steam.workerThread(wCtx)
	}
	return &steam, nil
}

// Stop steam seeker for all workers
func (steam *SteamSeeker) Stop() error {
	for i := 0; i < steam.config.WorkerNum; i++ {
		steam.workerQuit <- struct{}{}
	}
	return nil
}

// WaitUntilDone all steam seeker workers done their work
func (steam *SteamSeeker) WaitUntilDone() error {
	c := 0
	for {
		select {
		case <-steam.workerDone:
			c++
		}
		if c == steam.config.WorkerNum {
			return nil
		}
	}
}

func (steam *SteamSeeker) getSteamAppList(ctx context.Context) error {
	log.Printf("steam seeker getting app list")
	req, err := http.NewRequest("GET", steam.config.Portal+pathGetAppList, nil)
	if err != nil {
		return err
	}
	return httpDo(ctx, req, steam.client, steam.processSteamAppList)
}

func (steam *SteamSeeker) processSteamAppList(resp *http.Response, err error) error {
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

	oldList, err := steam.store.GetGameList("steam")
	if err != nil {
		return err
	}
	log.Printf("processSteamAppList: oldGameList len: %d, newGameList len: %d", len(oldList), len(gameList))
	diff, err := steam.createSeekerQueue(oldList, gameList)
	if err != nil {
		return err
	}
	if len(diff) > 0 {
		steam.store.SaveGameList("steam", gameList)
	}
	return nil
}

func (steam *SteamSeeker) createSeekerQueue(oldList map[int]string, newList map[int]string) (map[int]string, error) {
	ret := make(map[int]string)
	for k, v := range newList {
		if _, ok := oldList[k]; !ok {
			ret[k] = v
		}
	}
	go func() {
		for k := range ret {
			steam.queue <- k
		}
		for i := 0; i < steam.config.WorkerNum; i++ {
			steam.workerQuit <- struct{}{}
		}
		close(steam.queue)
	}()
	return ret, nil
}

func (steam *SteamSeeker) getSteamAppDetail(ctx context.Context, appid int) error {
	log.Printf("workerThread[%d] getting app detail: %d", ctx.Value(workerIDKey), appid)
	req, err := http.NewRequest("GET", pathGetAppDetail, nil)
	q := req.URL.Query()
	q.Add("appids", strconv.FormatInt(int64(appid), 10))
	req.URL.RawQuery = q.Encode()
	if err != nil {
		return err
	}
	return httpDo(ctx, req, steam.client, steam.processSteamAppDetail)
}

func (steam *SteamSeeker) processSteamAppDetail(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	res, err := steam.parseSteamAppDetail(body)
	if err != nil {
		return err
	}
	sad := res.Data
	var reqAge int64
	var err1 error
	switch v := sad.RequiredAge.(type) {
	case int:
		reqAge = int64(sad.RequiredAge.(int))
	case string:
		if reqAge, err1 = strconv.ParseInt(sad.RequiredAge.(string), 10, 32); err1 != nil {
			return err1
		}
	case float64:
		reqAge = int64(sad.RequiredAge.(float64))
	default:
		log.Printf("RequiredAge unknown type %T!", v)
		return errors.New("can't parse required age")
	}
	gr := store.GameRecord{
		Name:        sad.Name,
		RequiredAge: int(reqAge),
		Description: sad.DetailedDescription,
		About:       sad.AboutTheGame,
		Languages:   sad.SupportLanguages,
		Developers:  sad.Developers,
		Publishers:  sad.Publishers,
	}
	steam.workerReturn <- gr
	return nil
}

func (steam *SteamSeeker) parseSteamAppDetail(data []byte) (steamAppDetail, error) {
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
	return steamAppDetail{}, ErrSteamFailReponse
}

func (steam *SteamSeeker) workerThread(ctx context.Context) error {
	defer func() {
		steam.workerDone <- struct{}{}
	}()
	rc := steam.config.RetryCount
	for {
		select {
		case appID := <-steam.queue:
			for i := 0; ; i++ {
				if err := steam.getSteamAppDetail(ctx, appID); err == nil || (rc > 0 && i >= rc) {
					break
				} else {
					if err == ErrSteamRateLimit {
						i = 0
					}

					log.Printf("workerThread[%d] getSteamAppDetail err: %v appid: %d retry count: %d", ctx.Value(workerIDKey).(int), err, appID, i)
					time.Sleep(steam.config.RetryInterval)
				}
			}
		case <-steam.workerQuit:
			log.Printf("workerThread %d quiting", ctx.Value(workerIDKey).(int))
			return nil
		}
	}
}
