package main

import (
	// "./buildtree"
	// "./hanlp"
	// "./strsim"
	// "fmt"
	"./nlphelper"
)

func main() {
	// strMap, _, _, iNodeToWords, wordFreq, wordVector := buildtree.LoadTree("company_")
	// base := "北京弘高创意建筑设计股份有限公司"
	// baseSplit := iNodeToWords[strMap[base]]
	// str := "好莱客创意"
	// strSplit := hanlp.StrSegment(str)
	// fmt.Println(strsim.SplitVectorDis(strSplit, baseSplit, wordFreq, wordVector))
	// build("../all_companies.csv", nlphelper.NlpUsingHanlp)
	searchSplit("search.txt", nlphelper.NlpUsingHanlp)
}
