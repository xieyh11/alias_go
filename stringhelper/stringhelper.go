package stringhelper

import (
	"strconv"
	"strings"
)

var ChinesePunc = []string{"。", "，", "（", "）", "？", "！", "、", "；", "：", "”", "“", "‘", "’", "——", "……", "《", "》", "<", ">"}

func IsChinesePunc(str string) bool {
	for _, punc := range ChinesePunc {
		if punc == str {
			return true
		}
	}
	return false
}

var EnglishPunc = []string{".", ",", "(", ")", ":", ";", "<", ">"}

func IsEnglishPunc(str string) bool {
	for _, punc := range EnglishPunc {
		if punc == str {
			return true
		}
	}
	return false
}

func RemoveChinesePunc(word string) string {
	for _, punc := range ChinesePunc {
		if strings.Contains(word, punc) {
			word = strings.Replace(word, punc, "", -1)
		}
	}
	return word
}

func RemoveEnglishPunc(word string) string {
	for _, punc := range EnglishPunc {
		if strings.Contains(word, punc) {
			word = strings.Replace(word, punc, "", -1)
		}
	}
	return word
}

func RemovePunc(word string) string {
	word = RemoveChinesePunc(word)
	word = RemoveEnglishPunc(word)
	return word
}

func RemovePuncFromWords(words []string) []string {
	wordsLen := len(words)
	for i := 0; i < wordsLen; i++ {
		if IsChinesePunc(words[i]) || IsEnglishPunc(words[i]) {
			words = append(words[:i], words[i+1:]...)
			wordsLen--
		} else {
			words[i] = RemovePunc(words[i])
		}
	}
	return words
}

func ParseStringArray(str string, hasBracket bool, seg string) []string {
	if hasBracket {
		strRune := []rune(str)
		strRune = strRune[1 : len(strRune)-2]
		str = string(strRune)
	}
	return strings.Split(str, seg)
}

func ParseFloatArray(str string, hasBracket bool, seg string) []float64 {
	if hasBracket {
		strRune := []rune(str)
		strRune = strRune[1 : len(strRune)-2]
		str = string(strRune)
	}
	strSeg := strings.Split(str, seg)
	res := make([]float64, len(strSeg))
	for i := range strSeg {
		res[i], _ = strconv.ParseFloat(strSeg[i], 64)
	}
	return res
}
