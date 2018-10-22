package hanlp

import (
	"../../floatvector"
	"../../stringhelper"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type HanLPConfig struct {
	UrlConfig      HanLPUrlConfig      `json:"url_config"`
	ResponseConfig HanLPResponseConfig `json:"response_config"`
}

type HanLPUrlConfig struct {
	Host        string `json:"url_host"`
	Segment     string `json:"segment_url"`
	Word2Vector string `json:"word2vector_url"`
	SegmentKey  string `json:"segment_keyword"`
	Word2VecKey string `json:"word2vector_keyword"`
}

type HanLPResponseConfig struct {
	ArraySeg    string `json:"array_segment"`
	Word2VecDim int    `json:"word2vector_dim"`
}

func NewHanLPConfig(config string) *HanLPConfig {
	if len(config) == 0 {
		config = "nlphelper/hanlp/config.json"
	}

	file, err := os.Open(config)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	defer file.Close()
	jsonDec := json.NewDecoder(file)
	hanlpConfig := new(HanLPConfig)
	jsonDec.Decode(hanlpConfig)
	return hanlpConfig
}

func (hanlp *HanLPConfig) Segment(str string) []string {
	strEn, _ := url.ParseQuery(hanlp.UrlConfig.SegmentKey + "=" + str)
	segPrefix := hanlp.UrlConfig.Host + hanlp.UrlConfig.Segment + "?"
	segQuery := segPrefix + strEn.Encode()
	response, err := http.Get(segQuery)
	tryTimes := 10
	for tryTimes > 0 && (err != nil || response.StatusCode != http.StatusOK) {
		tryTimes--
		response, err = http.Get(segQuery)
		time.Sleep(time.Duration(rand.Int()%5) * time.Second)
	}
	if err != nil || response.StatusCode != http.StatusOK {
		panic(fmt.Sprintln(err))
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return stringhelper.ParseStringArray(string(body), true, hanlp.ResponseConfig.ArraySeg)
}

func (hanlp *HanLPConfig) ToVector(word string) []float64 {
	wordEn, _ := url.ParseQuery(hanlp.UrlConfig.Word2VecKey + "=" + word)
	wordPrefix := hanlp.UrlConfig.Host + hanlp.UrlConfig.Word2Vector + "?"
	wordQuery := wordPrefix + wordEn.Encode()
	response, err := http.Get(wordQuery)
	tryTimes := 10
	for tryTimes > 0 && (err != nil || response.StatusCode != http.StatusOK) {
		tryTimes--
		response, err = http.Get(wordQuery)
		time.Sleep(time.Duration(rand.Int()%5) * time.Second)
	}
	if err != nil || response.StatusCode != http.StatusOK {
		panic(fmt.Sprintln(err))
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	if strings.Contains(string(body), "null") {
		wordRune := []rune(word)
		if len(wordRune) == 1 {
			return make([]float64, hanlp.ResponseConfig.Word2VecDim)
		}
		vecRes := []float64{}
		for i := range wordRune {
			tempR := hanlp.ToVector(string(wordRune[i]))
			if len(vecRes) == 0 {
				vecRes = tempR
			} else {
				floatvector.AddVectorsInPlace(vecRes, tempR)
			}
		}
		return vecRes
	} else {
		return stringhelper.ParseFloatArray(string(body), true, hanlp.ResponseConfig.ArraySeg)
	}
}
