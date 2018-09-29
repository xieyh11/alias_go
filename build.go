package main

import (
	"fmt"
	"math"

	"./buildtree"
)

func build() {
	names := buildtree.ReadCsv("../all_companies.csv")
	wordFre, splitWords := buildtree.CountWordFrequence(names)
	strMap, iNodeMap, _, iNodeToWords := buildtree.BuildMapTree(names, splitWords)
	buildtree.DumpTree(strMap, iNodeMap, iNodeToWords, wordFre, "company_")

	strMap1, iNodeMap1, _, iNodeToWords1, wordFre1 := buildtree.LoadTree("company_")
	for k, v := range strMap {
		if v1, ok := strMap1[k]; !ok || (v != v1) {
			fmt.Println("String Map to INode error: " + k)
		}
	}
	for k, v := range iNodeMap {
		if v1, ok := iNodeMap1[k]; !ok {
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
	for k, v := range iNodeToWords {
		if v1, ok := iNodeToWords1[k]; ok {
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
	for k, v := range wordFre {
		if v1, ok := wordFre1[k]; ok {
			if math.Abs(v-v1) > 1e-5 {
				fmt.Println("Word's Freq is not same: ", k, v, v1)
			}
		} else {
			fmt.Println("Word's freq is not in the file: ", k)
		}
	}
}
