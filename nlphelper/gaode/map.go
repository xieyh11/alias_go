package gaode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/chromedp/chromedp"
)

type MapConfig struct {
	UrlConfig    MapUrlConfig       `json:"url_config"`
	chromeCtx    context.Context    `json:"-"`
	chromeIns    *chromedp.CDP      `json:"-"`
	chromeCancel context.CancelFunc `json:"-"`
}

type MapUrlConfig struct {
	Host   string `json:"url_host"`
	MapKey string `json:"map_keyword"`
}

func NewMapConfig(config string) *MapConfig {
	if len(config) == 0 {
		config = "nlphelper/gaode/config.json"
	}
	configFile, err := os.Open(config)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	defer configFile.Close()
	jsonDec := json.NewDecoder(configFile)
	res := new(MapConfig)
	jsonDec.Decode(res)
	res.chromeCtx, res.chromeCancel = context.WithCancel(context.Background())
	res.chromeIns, err = chromedp.New(res.chromeCtx)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	return res
}

type MapPoi struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	PName    string `json:"pname"`
	CityName string `json:"cityname"`
	AdName   string `json:"adname"`
}

func (mapConfig *MapConfig) Segment(str string) []string {
	strEn, _ := url.ParseQuery(mapConfig.UrlConfig.MapKey + "=" + str)
	addPrefix := mapConfig.UrlConfig.Host + "?"
	addQuery := addPrefix + strEn.Encode()
	var res string
	err := mapConfig.chromeIns.Run(mapConfig.chromeCtx, chromedp.Tasks{
		chromedp.Navigate(addQuery),
		chromedp.Text(`#search`, &res, chromedp.NodeVisible, chromedp.ByID),
	})
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	if len(res) == 0 {
		return []string{}
	} else if res == "null" {
		strRune := []rune(str)
		if len(strRune) < 2 {
			return []string{}
		} else {
			return mapConfig.Segment(string(strRune[:len(strRune)-1]))
		}
	} else {
		var mapRes MapPoi
		err = json.Unmarshal([]byte(res), &mapRes)
		if err != nil {
			panic(fmt.Sprintln(err))
		}
		return []string{mapRes.PName, mapRes.CityName, mapRes.AdName, mapRes.Address, mapRes.Name}
	}
}
