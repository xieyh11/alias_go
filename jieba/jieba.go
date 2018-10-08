package jieba

import (
	"bufio"
	"github.com/yanyiwu/gojieba"
	"os"
)

func NewJieba() *gojieba.Jieba {
	res := gojieba.NewJieba()

	userDict, err := os.Open("user.dict")
	defer userDict.Close()

	if err != nil {
		return res
	}

	userDictScan := bufio.NewScanner(userDict)
	userDictScan.Split(bufio.ScanLines)

	for userDictScan.Scan() {
		line := userDictScan.Text()
		res.AddWord(line)
	}
	return res
}
