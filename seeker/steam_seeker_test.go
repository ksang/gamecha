package seeker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ksang/gamecha/store"
)

func TestGetSteamAppList(t *testing.T) {
	s, _ := store.NewDummyStore(store.Config{})
	steam := &SteamConfig{
		Portal:    "http://api.steampowered.com",
		Key:       "",
		ThreadNum: 10,
		store:     s,
		queue:     make(chan int),
	}
	timeout, _ := time.ParseDuration("60s")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := steam.getSteamAppList(ctx); err != nil {
		t.Errorf("getSteamAppList err: %v", err)
	}
	fmt.Println("Games to be collected:")
	for id := range steam.queue {
		fmt.Printf("%d ", id)
	}
}

func TestParseSteamAppDetail(t *testing.T) {
	dataStr := `{
  "10": {
    "success": true,
    "data": {
      "type": "game",
      "name": "Counter-Strike",
      "steam_appid": 10,
      "required_age": 0,
      "is_free": false,
      "detailed_description": "Play the world's number 1 online action game. Engage in an incredibly realistic brand of terrorist warfare in this wildly popular team-based game. Ally with teammates to complete strategic missions. Take out enemy sites. Rescue hostages. Your role affects your team's success. Your team's success affects your role.",
      "about_the_game": "Play the world's number 1 online action game. Engage in an incredibly realistic brand of terrorist warfare in this wildly popular team-based game. Ally with teammates to complete strategic missions. Take out enemy sites. Rescue hostages. Your role affects your team's success. Your team's success affects your role.",
      "short_description": "Play the world's number 1 online action game. Engage in an incredibly realistic brand of terrorist warfare in this wildly popular team-based game. Ally with teammates to complete strategic missions. Take out enemy sites. Rescue hostages. Your role affects your team's success. Your team's success affects your role.",
      "supported_languages": "English<strong>*<\/strong>, French<strong>*<\/strong>, German<strong>*<\/strong>, Italian<strong>*<\/strong>, Spanish - Spain<strong>*<\/strong>, Simplified Chinese<strong>*<\/strong>, Traditional Chinese<strong>*<\/strong>, Korean<strong>*<\/strong><br><strong>*<\/strong>languages with full audio support",
      "header_image": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/header.jpg?t=1528733245",
      "website": null,
      "pc_requirements": {
        "minimum": "\r\n\t\t\t<p><strong>Minimum:<\/strong> 500 mhz processor, 96mb ram, 16mb video card, Windows XP, Mouse, Keyboard, Internet Connection<br \/><\/p>\r\n\t\t\t<p><strong>Recommended:<\/strong> 800 mhz processor, 128mb ram, 32mb+ video card, Windows XP, Mouse, Keyboard, Internet Connection<br \/><\/p>\r\n\t\t\t"
      },
      "mac_requirements": {
        "minimum": "Minimum: OS X  Snow Leopard 10.6.3, 1GB RAM, 4GB Hard Drive Space,NVIDIA GeForce 8 or higher, ATI X1600 or higher, or Intel HD 3000 or higher Mouse, Keyboard, Internet Connection"
      },
      "linux_requirements": {
        "minimum": "Minimum: Linux Ubuntu 12.04, Dual-core from Intel or AMD at 2.8 GHz, 1GB Memory, nVidia GeForce 8600\/9600GT, ATI\/AMD Radeaon HD2600\/3600 (Graphic Drivers: nVidia 310, AMD 12.11), OpenGL 2.1, 4GB Hard Drive Space, OpenAL Compatible Sound Card"
      },
      "developers": [
        "Valve"
      ],
      "publishers": [
        "Valve"
      ],
      "price_overview": {
        "currency": "CNY",
        "initial": 3700,
        "final": 3700,
        "discount_percent": 0,
        "initial_formatted": "",
        "final_formatted": "\u00a5 37"
      },
      "packages": [
        7
      ],
      "package_groups": [
        {
          "name": "default",
          "title": "Buy Counter-Strike",
          "description": "",
          "selection_text": "Select a purchase option",
          "save_text": "",
          "display_type": 0,
          "is_recurring_subscription": "false",
          "subs": [
            {
              "packageid": 7,
              "percent_savings_text": "",
              "percent_savings": 0,
              "option_text": "Counter-Strike: Condition Zero - \u00a5 37",
              "option_description": "",
              "can_get_free_license": "0",
              "is_free_license": false,
              "price_in_cents_with_discount": 3700
            }
          ]
        }
      ],
      "platforms": {
        "windows": true,
        "mac": true,
        "linux": true
      },
      "metacritic": {
        "score": 88,
        "url": "https:\/\/www.metacritic.com\/game\/pc\/counter-strike?ftag=MCD-06-10aaa1f"
      },
      "categories": [
        {
          "id": 1,
          "description": "Multi-player"
        },
        {
          "id": 36,
          "description": "Online Multi-Player"
        },
        {
          "id": 37,
          "description": "Local Multi-Player"
        },
        {
          "id": 8,
          "description": "Valve Anti-Cheat enabled"
        }
      ],
      "genres": [
        {
          "id": "1",
          "description": "Action"
        }
      ],
      "screenshots": [
        {
          "id": 0,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000132.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000132.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 1,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000133.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000133.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 2,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000134.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000134.1920x1080* Connection #0 to host store.steampowered.com left intact .jpg?t=1528733245"
        },
        {
          "id": 3,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000135.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000135.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 4,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000136.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000000136.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 5,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002540.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002540.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 6,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002539.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002539.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 7,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002538.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002538.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 8,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002537.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002537.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 9,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002536.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002536.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 10,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002541.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002541.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 11,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002542.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002542.1920x1080.jpg?t=1528733245"
        },
        {
          "id": 12,
          "path_thumbnail": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002543.600x338.jpg?t=1528733245",
          "path_full": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/0000002543.1920x1080.jpg?t=1528733245"
        }
      ],
      "recommendations": {
        "total": 65110
      },
      "achievements": {
        "total": 0
      },
      "release_date": {
        "coming_soon": false,
        "date": "1 Nov, 2000"
      },
      "support_info": {
        "url": "http:\/\/steamcommunity.com\/app\/10",
        "email": ""
      },
      "background": "https:\/\/steamcdn-a.akamaihd.net\/steam\/apps\/10\/page_bg_generated_v6b.jpg?t=1528733245",
      "content_descriptors": {
        "ids": [
          2,
          5
        ],
        "notes": "Includes intense violence and blood."
      }
    }
  }
}`
	steam := &SteamConfig{}
	appDetail, err := steam.parseSteamAppDetail([]byte(dataStr))
	if err != nil {
		t.Errorf("TestParseSteamAppDetail err: %v", err)
	}
	t.Logf("TestParseSteamAppDetail result: %#v", appDetail)
}

func TestGetSteamAppDetail(t *testing.T) {
	steam := &SteamConfig{
		Portal:       "http://api.steampowered.com",
		Key:          "",
		queue:        make(chan int),
		workerReturn: make(chan store.GameRecord, 10),
	}
	var tests = []struct {
		id int
	}{
		{
			10,
		},
		{
			1050240,
		},
	}

	timeout, _ := time.ParseDuration("20s")
	for caseid, c := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := steam.getSteamAppDetail(ctx, c.id); err != nil {
			t.Errorf("case: #%d getSteamAppDetail err: %v", caseid+1, err)
			continue
		}
		gr := <-steam.workerReturn
		t.Logf("case: #%d got game record: %#v", caseid+1, gr)
	}

}
