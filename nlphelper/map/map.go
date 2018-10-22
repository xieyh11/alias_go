package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

type MapConfig struct {
	UrlConfig MapUrlConfig `json:"url_config"`
}

type MapUrlConfig struct {
	Host   string `json:"url_host"`
	MapKey string `json:"map_keyword"`
}

func NewMapConfig(config string) *MapConfig {
	if len(config) == 0 {
		config = "nlphelper/map/config.json"
	}
	configFile, err := os.Open(config)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	defer configFile.Close()
	jsonDec := json.NewDecoder(configFile)
	res := new(MapConfig)
	jsonDec.Decode(res)
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
}
