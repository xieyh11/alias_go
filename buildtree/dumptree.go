package buildtree

import (
	"../floatvector"
	"os"
	"strconv"
	"strings"
)

const StrToINodeSuffix = "str_map.pack"
const MapTreeSuffix = "tree.pack"
const INodeToWordsSuffix = "words.pack"
const WordsWeightSuffix = "weight.pack"
const WordsVectorSuffix = "vector.pack"

func DumpTree(mapTree MapTree, filePrefix string) {
	strMapFile, _ := os.Create(filePrefix + StrToINodeSuffix)
	treeFile, _ := os.Create(filePrefix + MapTreeSuffix)

	defer strMapFile.Close()
	defer treeFile.Close()

	for k, v := range mapTree.StrToINode {
		tempTreeNode := mapTree.INodeToTreeNode[v]
		strMapFile.WriteString(k + " " + strconv.FormatUint(v, 10) + "\n")
		strNode := strconv.FormatUint(tempTreeNode.INode, 10)
		if tempTreeNode.Parent != nil {
			strNode += " " + strconv.FormatUint(tempTreeNode.Parent.INode, 10)
		}
		treeFile.WriteString(strNode + "\n")
	}

}

func DumpMessage(strMessage StrMessage, filePrefix string) {
	wordsFile, _ := os.Create(filePrefix + INodeToWordsSuffix)
	freqFile, _ := os.Create(filePrefix + WordsWeightSuffix)
	vecFile, _ := os.Create(filePrefix + WordsVectorSuffix)

	defer wordsFile.Close()
	defer freqFile.Close()
	defer vecFile.Close()

	for k, v := range strMessage.INodeToWords {
		wordsFile.WriteString(strconv.FormatUint(k, 10) + " " + strings.Join(v, " ") + "\n")
	}
	for k, v := range strMessage.WordsWeight {
		freqFile.WriteString(k + " " + strconv.FormatFloat(v, 'f', -1, 64) + "\n")
	}
	for k, v := range strMessage.WordsVector {
		vecFile.WriteString(k + " " + floatvector.VectorToString(v) + "\n")
	}
}
