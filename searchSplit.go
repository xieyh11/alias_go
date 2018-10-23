package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"./buildtree"
	"./nlphelper"
	"./nlphelper/gaode"
	"./nlphelper/hanlp"
	"./nlphelper/jieba"
	"./searchtree"
)

func searchSplit(file, filePrefix string, nlpUsing int) {
	mapTree := buildtree.LoadTree(filePrefix)
	strMessage := buildtree.LoadMessage(filePrefix)
	indexStr := make([]string, len(mapTree.StrToINode))
	index := 0
	for k, _ := range mapTree.StrToINode {
		indexStr[index] = k
		index++
	}
	searchFile, _ := os.Open(file)
	defer searchFile.Close()

	scanner := bufio.NewScanner(searchFile)
	scanner.Split(bufio.ScanLines)

	//res := make([][]string, 0)
	// jieba := jieba.NewJieba()
	// defer jieba.Free()
	var segmentor nlphelper.WordsSegment
	var word2vector nlphelper.Word2Vector
	hanlpConfig := hanlp.NewHanLPConfig("")
	word2vector = hanlpConfig
	switch nlpUsing {
	case nlphelper.NlpUsingJieba:
		segmentor = jieba.NewJieba()
	case nlphelper.NlpUsingGaode:
		segmentor = gaode.NewMapConfig("")
	default:
		segmentor = hanlpConfig
	}

	for scanner.Scan() {
		line := scanner.Text()
		name := segmentor.Segment(line)
		nameVector := make([][]float64, len(name))
		for i := range name {
			nameVector[i] = word2vector.ToVector(name[i])
		}
		scores := searchtree.SearchSplitStrings(indexStr, mapTree, strMessage, name, nameVector, 10)
		topK := searchtree.TopKScores(scores, 20)
		for i := range topK {
			str := indexStr[topK[i]] + " " + strconv.FormatFloat(scores[topK[i]], 'f', -1, 64) + " "
			iNode := mapTree.StrToINode[indexStr[topK[i]]]
			iTreeNode := mapTree.INodeToTreeNode[iNode]
			rootOf := iTreeNode.ToRoot()
			str += mapTree.INodeToStr[rootOf.INode]
			fmt.Println(str)
		}
		fmt.Println()
	}
}
