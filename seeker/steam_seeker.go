package seeker

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	// ErrSteamQuitTimeout indicates steam seeker took too long to graceful quit
	ErrSteamQuitTimeout = errors.New("steam seeker grace quit timed out")
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
	errc         chan error
	workerDone   chan struct{}
	workerReturn chan store.GameRecord
	debugLog     *log.Logger
	infoLog      *log.Logger
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
		errc:         make(chan error),
		store:        db,
		workerReturn: make(chan store.GameRecord, cfg.WorkerNum),
		workerDone:   make(chan struct{}, cfg.WorkerNum),
		debugLog:     log.New(os.Stdout, "SteamSeeker DEBUG:", log.LstdFlags|log.Lshortfile),
		infoLog:      log.New(os.Stdout, "SteamSeeker INFO:", log.LstdFlags|log.Lshortfile),
	}
	if err := steam.getSteamAppList(ctx); err != nil {
		return nil, err
	}
	go func() {
		steam.errc <- steam.storeRecord(ctx)
	}()
	for i := 0; i < cfg.WorkerNum; i++ {
		wCtx := context.WithValue(ctx, workerIDKey, i)
		go steam.workerThread(wCtx)
	}

	return &steam, nil
}

// WaitUntilDone all steam seeker workers done their work
func (steam *SteamSeeker) WaitUntilDone(ctx context.Context) error {
	allWorkerDone := make(chan struct{})
	go func() {
		c := 0
		for {
			select {
			case <-steam.workerDone:
				c++
			}
			if c == steam.config.WorkerNum {
				allWorkerDone <- struct{}{}
				return
			}
		}
	}()
	for {
		select {
		case <-allWorkerDone:
			return nil
		case <-ctx.Done():
			select {
			case <-time.After(3000000000):
				return ErrSteamQuitTimeout
			case <-allWorkerDone:
				return nil
			}
		case err := <-steam.errc:
			return err
		}
	}
}

func (steam *SteamSeeker) getSteamAppList(ctx context.Context) error {
	steam.debugLog.Printf("getting app list")
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

	oldList, err := steam.store.GetSavedGameList("steam")
	if err != nil {
		return err
	}
	steam.debugLog.Printf("processSteamAppList: oldGameList len: %d, newGameList len: %d", len(oldList), len(gameList))
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
		close(steam.queue)
	}()
	return ret, nil
}

func (steam *SteamSeeker) getSteamAppDetail(ctx context.Context, appid int) error {
	steam.debugLog.Printf("workerThread[%d] getting app detail: %d", ctx.Value(workerIDKey), appid)
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
		steam.infoLog.Printf("RequiredAge unknown type %T!", v)
		return errors.New("can't parse required age")
	}
	gr := store.GameRecord{
		Name:        sad.Name,
		ID:          sad.Appid,
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
	id := ctx.Value(workerIDKey).(int)
	for {
		select {
		case appID := <-steam.queue:
			// Retrying get one app detail several times according to config
			for i := 0; ; i++ {
				if err := steam.getSteamAppDetail(ctx, appID); err == nil || (rc > 0 && i >= rc) {
					break
				} else {
					// If err is due to steam api rate limit, keep retry
					if err == ErrSteamRateLimit {
						i = 0
					}
					steam.debugLog.Printf("workerThread[%d] getSteamAppDetail err: %v appid: %d count: %d", id, err, appID, i)
					select {
					case <-ctx.Done():
						steam.infoLog.Printf("workerThread[%d] signaled to quit", id)
						return nil
					case <-time.After(steam.config.RetryInterval):
						continue
					}
				}
			}
		case <-ctx.Done():
			steam.infoLog.Printf("workerThread[%d] signaled to quit", id)
			return nil
		}
	}
}

func (steam *SteamSeeker) storeRecord(ctx context.Context) error {
	for {
		select {
		case gr := <-steam.workerReturn:
			if err := steam.store.SaveGameRecord("steam", strconv.FormatInt(int64(gr.ID), 10), gr); err != nil {
				steam.infoLog.Printf("failed to save game record, appid: %d", gr.ID)
				return err
			}

		case <-ctx.Done():
			steam.infoLog.Printf("store record process signaled to quit")
			return nil
		}
	}
}
