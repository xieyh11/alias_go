package jieba

import (
	"bufio"
	"github.com/yanyiwu/gojieba"
	"os"
)

type Jieba gojieba.Jieba

func NewJieba() *Jieba {
	res := gojieba.NewJieba()

	userDict, err := os.Open("user.dict")
	defer userDict.Close()

	if err != nil {
		return (*Jieba)(res)
	}

	userDictScan := bufio.NewScanner(userDict)
	userDictScan.Split(bufio.ScanLines)

	for userDictScan.Scan() {
		line := userDictScan.Text()
		res.AddWord(line)
	}
	return (*Jieba)(res)
}

func (a *Jieba) Segment(str string) []string {
	return (*gojieba.Jieba)(a).Cut(str, true)
}
