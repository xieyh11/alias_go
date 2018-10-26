package main

import "./strsim"

func main() {
	// filePrefix := "address_gaode_"
	// mapTree := buildtree.LoadTree(filePrefix)
	// strMessage := buildtree.LoadMessage(filePrefix)
	// address := "武汉市武昌区友谊大道特1号广达大厦写字楼"
	// compare := "香港皇后大道东183号合和中心64楼"
	// mapConfig := gaode.NewMapConfig("")
	// words := mapConfig.Segment(address)
	// fmt.Println(strsim.MapSegmentSimiliarity(strMessage.INodeToWords[mapTree.StrToINode[compare]], words, strsim.MapSimiliarityLevelThree))
	// update("../address.csv", "address_gaode_", nlphelper.NlpUsingGaode)
	searchMap("searchMap.txt", "address_gaode_", strsim.MapSimiliarityLevelThree)
}
