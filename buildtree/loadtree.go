package buildtree

import (
	"../maptree"
	"bufio"
	"os"
	"strconv"
	"strings"
)

func LoadTree(filePrefix string) MapTree {
	mapTree := MapTree{INodeToStr: make(map[uint64]string), StrToINode: make(map[string](uint64)), INodeToTreeNode: make(map[uint64](*maptree.MapTreeNode))}

	strMapFile, _ := os.Open(filePrefix + StrToINodeSuffix)
	defer strMapFile.Close()

	strScanner := bufio.NewScanner(strMapFile)
	strScanner.Split(bufio.ScanLines)

	for strScanner.Scan() {
		line := strScanner.Text()
		if len(line) > 0 {
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				iNode, _ := strconv.ParseUint(parts[1], 10, 64)
				mapTree.StrToINode[parts[0]] = iNode
				mapTree.INodeToStr[iNode] = parts[0]
			}
		}
	}

	treeFile, _ := os.Open(filePrefix + "tree.pack")
	defer treeFile.Close()

	treeScanner := bufio.NewScanner(treeFile)
	treeScanner.Split(bufio.ScanLines)

	for treeScanner.Scan() {
		line := treeScanner.Text()
		if len(line) > 0 {
			parts := strings.Split(line, " ")
			parentNode := (*maptree.MapTreeNode)(nil)
			if len(parts) == 2 {
				iNode, _ := strconv.ParseUint(parts[1], 10, 64)
				if _, ok := mapTree.INodeToTreeNode[iNode]; ok {
					parentNode = mapTree.INodeToTreeNode[iNode]
				} else {
					parentNode = maptree.NewMapTreeNode(iNode)
					mapTree.INodeToTreeNode[iNode] = parentNode
				}
			}
			iNode, _ := strconv.ParseUint(parts[0], 10, 64)
			if _, ok := mapTree.INodeToTreeNode[iNode]; ok {
				mapTree.INodeToTreeNode[iNode].Parent = parentNode
			} else {
				newNode := maptree.NewMapTreeNodeWithParent(iNode, parentNode)
				mapTree.INodeToTreeNode[iNode] = newNode
			}
		}
	}

	return mapTree
}

func LoadMessage(filePrefix string) StrMessage {
	strMessage := StrMessage{INodeToWords: make(map[uint64]([]string)), WordsVector: make(map[string][]float64), WordsWeight: make(map[string]float64)}
	wordsFile, _ := os.Open(filePrefix + INodeToWordsSuffix)
	defer wordsFile.Close()

	wordsScanner := bufio.NewScanner(wordsFile)
	wordsScanner.Split(bufio.ScanLines)

	for wordsScanner.Scan() {
		line := wordsScanner.Text()
		if len(line) > 0 {
			words := strings.Split(line, " ")
			iNode, _ := strconv.ParseUint(words[0], 10, 64)
			strMessage.INodeToWords[iNode] = words[1:]
		}
	}

	weightFile, _ := os.Open(filePrefix + WordsWeightSuffix)
	defer weightFile.Close()

	weightScanner := bufio.NewScanner(weightFile)
	weightScanner.Split(bufio.ScanLines)

	for weightScanner.Scan() {
		line := weightScanner.Text()
		if len(line) > 0 {
			splits := strings.Split(line, " ")
			weight, _ := strconv.ParseFloat(splits[1], 64)
			strMessage.WordsWeight[splits[0]] = weight
		}
	}

	vecFile, _ := os.Open(filePrefix + WordsVectorSuffix)
	defer vecFile.Close()

	vecScanner := bufio.NewScanner(vecFile)
	vecScanner.Split(bufio.ScanLines)

	for vecScanner.Scan() {
		line := vecScanner.Text()
		if len(line) > 0 {
			splits := strings.Split(line, " ")
			vector := make([]float64, len(splits)-1)
			for i := 1; i < len(splits); i++ {
				vector[i-1], _ = strconv.ParseFloat(splits[i], 64)
			}
			strMessage.WordsVector[splits[0]] = vector
		}
	}
	return strMessage
}
