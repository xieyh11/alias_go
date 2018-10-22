package main

import (
	"fmt"
	"math"

	"./buildtree"
)

func build(csvFile string, filePrefix string, nlpUsing int) {
	names := buildtree.ReadCsv(csvFile, true, true)
	mapTree := buildtree.BuildMapTree(names)
	strMessage := buildtree.EvaluateStrMessage(names, mapTree, nlpUsing)
	buildtree.DumpTree(mapTree, filePrefix)
	buildtree.DumpMessage(strMessage, filePrefix)

	mapTree1 := buildtree.LoadTree(filePrefix)
	strMessage1 := buildtree.LoadMessage(filePrefix)
	for k, v := range mapTree.StrToINode {
		if v1, ok := mapTree1.StrToINode[k]; !ok || (v != v1) {
			fmt.Println("String Map to INode error: " + k)
		}
	}
	for k, v := range mapTree.INodeToTreeNode {
		if v1, ok := mapTree1.INodeToTreeNode[k]; !ok {
			fmt.Println("INode doesn't appear in the tree: ", k)
		} else {
			if v.Parent == nil {
				if v1.Parent != nil {
					fmt.Println("INode is root in one tree but not in other three: ", k)
				}
			} else {
				if v1.Parent == nil {
					fmt.Println("INode is root in one tree but not in other three: ", k)
				} else {
					if v.Parent.INode != v1.Parent.INode {
						fmt.Println("INode's parent are same: ", k, v.Parent.INode, v1.Parent.INode)
					}
				}
			}
		}
	}
	for k, v := range strMessage.INodeToWords {
		if v1, ok := strMessage1.INodeToWords[k]; ok {
			equalSplit := true
			for i := range v {
				if v[i] != v1[i] {
					equalSplit = false
					break
				}
			}
			if !equalSplit {
				fmt.Println("Split is not equal: ", k, v, v1)
			}
		} else {
			fmt.Println("Words Split doesn't exist in the file: ", k)
		}
	}
	for k, v := range strMessage.WordsWeight {
		if v1, ok := strMessage1.WordsWeight[k]; ok {
			if math.Abs(v-v1) > 1e-5 {
				fmt.Println("Word's Freq is not same: ", k, v, v1)
			}
		} else {
			fmt.Println("Word's freq is not in the file: ", k)
		}
	}
	for k, v := range strMessage.WordsVector {
		if v1, ok := strMessage1.WordsVector[k]; ok {
			similiar := true
			if len(v) == len(v1) {
				for i := range v {
					if math.Abs(v[i]-v1[i]) > 1e-5 {
						similiar = false
					}
					if !similiar {
						break
					}
				}
			} else {
				similiar = false
			}
			if !similiar {
				fmt.Println("Word's Vector is not same: ", k, v, v1)
			}
		} else {
			fmt.Println("Word's Vector is not in the file: ", k)
		}
	}
}
