package hanlp

import (
	"../floatvector"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var urlPrefix = "http://180.8.50.64:8080/webhanlp/"
var segment = "segment"
var word2Vector = "word"
var segKey = "str"
var wordKey = "word"
var arraySeg = ", "
var vectorDim = 300

func parseArrayString(str string) []string {
	strRune := []rune(str)
	strRune = strRune[1 : len(strRune)-2]
	str = string(strRune)
	return strings.Split(str, arraySeg)
}

func parseArrayFloat(str string) []float64 {
	strRune := []rune(str)
	strRune = strRune[1 : len(strRune)-2]
	str = string(strRune)
	strSep := strings.Split(str, arraySeg)
	res := make([]float64, len(strSep))
	for i := range strSep {
		res[i], _ = strconv.ParseFloat(strSep[i], 64)
	}
	return res
}

func StrSegment(str string) []string {
	strEn, _ := url.ParseQuery(segKey + "=" + str)
	response, err := http.Get(urlPrefix + segment + "?" + strEn.Encode())
	tryTimes := 10
	for tryTimes > 0 && (err != nil || response.StatusCode != http.StatusOK) {
		tryTimes--
		response, err = http.Get(urlPrefix + segment + "?" + strEn.Encode())
		time.Sleep(5 * time.Second)
	}
	if err != nil || response.StatusCode != http.StatusOK {
		panic(fmt.Sprintln(err))
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return parseArrayString(string(body))
}

func Word2Vector(word string) []float64 {
	wordEn, _ := url.ParseQuery(wordKey + "=" + word)
	response, err := http.Get(urlPrefix + word2Vector + "?" + wordEn.Encode())
	tryTimes := 10
	for tryTimes > 0 && (err != nil || response.StatusCode != http.StatusOK) {
		tryTimes--
		response, err = http.Get(urlPrefix + word2Vector + "?" + wordEn.Encode())
		time.Sleep(5 * time.Second)
	}
	if err != nil || response.StatusCode != http.StatusOK {
		panic(fmt.Sprintln(err))
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	if strings.Contains(string(body), "null") {
		wordRune := []rune(word)
		if len(wordRune) == 1 {
			return make([]float64, vectorDim)
		}
		vecRes := []float64{}
		for i := range wordRune {
			tempR := Word2Vector(string(wordRune[i]))
			if len(vecRes) == 0 {
				vecRes = tempR
			} else {
				floatvector.AddVectorsInPlace(vecRes, tempR)
			}
		}
		return vecRes
	} else {
		return parseArrayFloat(string(body))
	}
}
