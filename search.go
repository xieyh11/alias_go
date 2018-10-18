package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"./buildtree"
	"./searchtree"
)

func search(file string) {
	filePrefix := "company_"
	mapTree := buildtree.LoadTree(filePrefix)
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

	for scanner.Scan() {
		line := scanner.Text()
		scores := searchtree.SearchStrings(indexStr, line, 10)
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
