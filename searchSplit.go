package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"./jieba"

	"./buildtree"
	"./searchtree"
)

func searchSplit() {
	strMap, iNodeMap, iNodeReverse, iNodeToWords, wordFreq := buildtree.LoadTree("company_")
	indexStr := make([]string, len(strMap))
	index := 0
	for k, _ := range strMap {
		indexStr[index] = k
		index++
	}
	searchFile, _ := os.Open("search.txt")
	defer searchFile.Close()

	scanner := bufio.NewScanner(searchFile)
	scanner.Split(bufio.ScanLines)

	//res := make([][]string, 0)
	jieba := jieba.NewJieba()
	defer jieba.Free()

	for scanner.Scan() {
		line := scanner.Text()
		name := jieba.Cut(line, true)
		scores := searchtree.SearchSplitStrings(indexStr, strMap, iNodeToWords, wordFreq, name, 10)
		topK := searchtree.TopKScores(scores, 20)
		for i := range topK {
			str := indexStr[topK[i]] + " " + strconv.FormatFloat(scores[topK[i]], 'f', -1, 64) + " "
			iNode := strMap[indexStr[topK[i]]]
			iTreeNode := iNodeMap[iNode]
			rootOf := iTreeNode.ToRoot()
			str += iNodeReverse[rootOf.INode]
			fmt.Println(str)
		}
		fmt.Println()
	}
}
