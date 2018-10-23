package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"./buildtree"
	"./nlphelper"
	"./nlphelper/gaode"
	"./searchtree"
)

func searchMap(file, filePrefix string, level int) {
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
	segmentor = gaode.NewMapConfig("")

	for scanner.Scan() {
		line := scanner.Text()
		name := segmentor.Segment(line)
		scores := searchtree.SearchMap(indexStr, mapTree, strMessage, name, level, 1)
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
