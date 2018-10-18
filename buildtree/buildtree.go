package buildtree

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"

	"../nlphelper"
	"../nlphelper/hanlp"
	"../nlphelper/jieba"

	"../maptree"
	"../stringhelper"
)

// Read a csv and get 2D array of company names with first column its simple name and second column its full name
func ReadCsv(filename string, hasIndex, hasColumnName bool) [][]string {
	csvFile, err := os.Open(filename)
	defer csvFile.Close()
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	err = nil
	res, err := csvReader.ReadAll()
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	if hasColumnName {
		res = res[1:]
	}
	if hasIndex {
		for i := range res {
			res[i] = res[i][1:]
		}
	}
	return res
}

func SplitStringIntoWords(csvDatas [][]string, segmenter nlphelper.WordsSegment) [][][]string {
	if len(csvDatas) == 0 || len(csvDatas[0]) == 0 {
		return [][][]string{}
	}
	rows, cols := len(csvDatas), len(csvDatas[0])
	res := make([][][]string, rows)
	for i := range res {
		res[i] = make([][]string, cols)
		for j := range res[i] {
			if csvDatas[i][j] == "" {
				res[i][j] = []string{}
			} else {
				res[i][j] = stringhelper.RemovePuncFromWords(segmenter.Segment(csvDatas[i][j]))
			}
		}
	}
	return res
}

const (
	WordsWeightEachRowOnce = iota
	WordsWeightEachStrOnce
	WordsWeightEachWord
)

func countWordsOnce(words []string, wordSet map[string]float64) {
	tempWord := make(map[string]bool)
	for _, word := range words {
		tempWord[word] = true
	}
	for k, _ := range tempWord {
		wordSet[k]++
	}
}

func EvaluateWordsWeight(words [][][]string, method int) map[string]float64 {
	res := make(map[string]float64)
	sum := 0.0
	switch method {
	case WordsWeightEachRowOnce:
		sum = float64(len(words))
		for _, row := range words {
			temp := make([]string, 0)
			for _, col := range row {
				temp = append(temp, col...)
			}
			if len(temp) == 0 {
				sum--
			} else {
				countWordsOnce(temp, res)
			}
		}
	default:
	}
	for k, v := range res {
		res[k] = 1 - v/sum
	}
	return res
}

type MapTree struct {
	StrToINode      map[string]uint64
	INodeToTreeNode map[uint64](*maptree.MapTreeNode)
	INodeToStr      map[uint64]string
}

type StrMessage struct {
	INodeToWords map[uint64][]string
	WordsVector  map[string][]float64
	WordsWeight  map[string]float64
}

func BuildMapTree(csvDatas [][]string) (mapTree MapTree) {
	mapTree.INodeToTreeNode = make(map[uint64](*maptree.MapTreeNode))
	mapTree.StrToINode = make(map[string](uint64))
	mapTree.INodeToStr = make(map[uint64]string)

	count := uint64(0)

	cols := len(csvDatas[0])
	for _, row := range csvDatas {
		if cols == 2 {
			if _, ok1 := mapTree.StrToINode[row[cols-1]]; ok1 {
				if _, ok2 := mapTree.StrToINode[row[0]]; ok2 {
					tempNode := mapTree.INodeToTreeNode[mapTree.StrToINode[row[0]]]
					tempNode = tempNode.ToRoot()
					tempName := mapTree.INodeToStr[tempNode.INode]
					if tempName != row[1] {
						fmt.Println("Error, simple name appear twice", row[0])
					}
				} else {
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

func EvaluateStrMessage(csvDatas [][]string, mapTree MapTree, nlpUsing int) (strMessage StrMessage) {
	strMessage.INodeToWords = make(map[uint64][]string)
	strMessage.WordsVector = make(map[string][]float64)

	var segmenter nlphelper.WordsSegment
	var word2vector nlphelper.Word2Vector
	switch nlpUsing {
	case nlphelper.NlpUsingJieba:
		segmenter = jieba.NewJieba()
		word2vector = hanlp.NewHanLPConfig("")
	case nlphelper.NlpUsingHanlp:
		fallthrough
	default:
		hanlpConfig := hanlp.NewHanLPConfig("")
		segmenter = hanlpConfig
		word2vector = hanlpConfig
	}

	splitStrings := SplitStringIntoWords(csvDatas, segmenter)
	strMessage.WordsWeight = EvaluateWordsWeight(splitStrings, WordsWeightEachRowOnce)

	for i := range csvDatas {
		for j := range csvDatas[i] {
			strMessage.INodeToWords[mapTree.StrToINode[csvDatas[i][j]]] = splitStrings[i][j]
		}
	}

	for k, _ := range strMessage.WordsWeight {
		strMessage.WordsVector[k] = word2vector.ToVector(k)
	}

	return
}
