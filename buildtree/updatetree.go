package buildtree

import (
	"fmt"
	"log"
	"os"

	"../maptree"
	"../nlphelper"
	"../nlphelper/gaode"
	"../nlphelper/hanlp"
	"../nlphelper/jieba"
)

func UpdateMapTree(csvDatas [][]string, mapTree MapTree) {
	count := uint64(0)
	for _, v := range mapTree.StrToINode {
		if count < v {
			count = v
		}
	}

	count++

	cols := len(csvDatas[0])
	for _, row := range csvDatas {
		if cols == 2 {
			if _, ok1 := mapTree.StrToINode[row[cols-1]]; ok1 {
				if _, ok2 := mapTree.StrToINode[row[0]]; !ok2 {
					parentNode := mapTree.INodeToTreeNode[mapTree.StrToINode[row[1]]]
					parentNode = parentNode.ToRoot()
					newNode := maptree.NewMapTreeNodeWithParent(count, parentNode)
					count++
					mapTree.StrToINode[row[0]] = newNode.INode
					mapTree.INodeToStr[newNode.INode] = row[0]
					mapTree.INodeToTreeNode[newNode.INode] = newNode
				}
			} else {
				if _, ok2 := mapTree.StrToINode[row[0]]; ok2 {
					fmt.Println("Error, simple name appear twice", row[0])
				} else {
					mapTree.StrToINode[row[1]] = count
					count++
					parentNode := maptree.NewMapTreeNode(mapTree.StrToINode[row[1]])
					newNode := maptree.NewMapTreeNodeWithParent(count, parentNode)
					count++

					mapTree.INodeToStr[parentNode.INode] = row[1]
					mapTree.INodeToStr[newNode.INode] = row[0]
					mapTree.StrToINode[row[0]] = newNode.INode
					mapTree.INodeToTreeNode[parentNode.INode] = parentNode
					mapTree.INodeToTreeNode[newNode.INode] = newNode
				}
			}
		} else {
			if _, ok := mapTree.StrToINode[row[0]]; !ok {
				newNode := maptree.NewMapTreeNode(count)
				count++

				mapTree.StrToINode[row[0]] = newNode.INode
				mapTree.INodeToTreeNode[newNode.INode] = newNode
				mapTree.INodeToStr[newNode.INode] = row[0]
			}
		}
	}

	return
}

func UpdateStrMessage(csvDatas [][]string, mapTree MapTree, strMessage StrMessage, nlpUsing int) {

	var segmenter nlphelper.WordsSegment
	var word2vector nlphelper.Word2Vector
	switch nlpUsing {
	case nlphelper.NlpUsingJieba:
		segmenter = jieba.NewJieba()
		word2vector = hanlp.NewHanLPConfig("")
	case nlphelper.NlpUsingGaode:
		segmenter = gaode.NewMapConfig("")
		word2vector = hanlp.NewHanLPConfig("")
	case nlphelper.NlpUsingHanlp:
		fallthrough
	default:
		hanlpConfig := hanlp.NewHanLPConfig("")
		segmenter = hanlpConfig
		word2vector = hanlpConfig
	}

	newCsvDatas := make([][]string, 0)
	for _, row := range csvDatas {
		if _, ok := strMessage.INodeToWords[mapTree.StrToINode[row[0]]]; !ok {
			newCsvDatas = append(newCsvDatas, row)
		} else {
			if len(row) == 2 {
				if _, ok := strMessage.INodeToWords[mapTree.StrToINode[row[1]]]; !ok {
					newCsvDatas = append(newCsvDatas, row)
				}
			}
		}
	}

	splitStrings := SplitStringIntoWords(newCsvDatas, segmenter)
	tempWeights := EvaluateWordsWeight(splitStrings, WordsWeightEachRowOnce)
	alreadyHave := float64(float64(len(csvDatas) - len(newCsvDatas)))
	newAdd := float64(len(newCsvDatas))

	for k, v := range tempWeights {
		if currentWeight, ok := strMessage.WordsWeight[k]; ok {
			strMessage.WordsWeight[k] = (v*alreadyHave + currentWeight*newAdd) / (alreadyHave + newAdd)
		} else {
			strMessage.WordsWeight[k] = v * alreadyHave / (alreadyHave + newAdd)
		}
	}

	logFile, _ := os.Create("log.pack")
	defer logFile.Close()
	log.SetOutput(logFile)

	for i := range newCsvDatas {
		for j := range newCsvDatas[i] {
			if len(splitStrings[i][j]) > 0 {
				strMessage.INodeToWords[mapTree.StrToINode[newCsvDatas[i][j]]] = splitStrings[i][j]
			} else {
				log.Print(newCsvDatas[i][j])
			}
		}
	}

	for k, _ := range strMessage.WordsWeight {
		if _, ok := strMessage.WordsVector[k]; !ok {
			strMessage.WordsVector[k] = word2vector.ToVector(k)
		}
	}

	return
}
