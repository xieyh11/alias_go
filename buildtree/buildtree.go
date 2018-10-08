package buildtree

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"../jieba"

	"../maptree"
)

// Read a csv and get 2D array of company names with first column its simple name and second column its full name
func ReadCsv(filename string) [][]string {
	csvFile, err := os.Open(filename)
	defer csvFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	err = nil
	res, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return res
	}
	if res[0][1] == "公司简称" {
		res = res[1:]
	}
	if len(res[0]) == 3 {
		for i := range res {
			res[i] = res[i][1:]
		}
	}
	return res
}

var ChinesePunc = []string{"。", "，", "（", "）", "？", "！", "、", "；", "：", "”", "“", "‘", "’", "——", "……", "《", "》", "<", ">"}

func isChinesePunc(str string) bool {
	for _, punc := range ChinesePunc {
		if punc == str {
			return true
		}
	}
	return false
}

var EnglishPunc = []string{".", ",", "(", ")", ":", ";", "<", ">"}

func isEnglishPunc(str string) bool {
	for _, punc := range EnglishPunc {
		if punc == str {
			return true
		}
	}
	return false
}

func CountWordFrequence(csvNames [][]string) (map[string]float64, [][]string) {
	wordSet := make(map[string]float64)
	splitWords := make([][]string, 0)

	jieba := jieba.NewJieba()
	defer jieba.Free()

	for _, row := range csvNames {
		words := jieba.Cut(row[1], true)
		for i := 0; i < len(words); {
			if isChinesePunc(words[i]) || isEnglishPunc(words[i]) {
				words = append(words[:i], words[i+1:]...)
			} else {
				i++
			}
		}
		tempWord := make(map[string]bool)
		for _, word := range words {
			tempWord[word] = true
		}
		for k, _ := range tempWord {
			if _, ok := wordSet[k]; ok {
				wordSet[k]++
			} else {
				wordSet[k] = 1
			}
		}
		splitWords = append(splitWords, words)
	}

	for k, v := range wordSet {
		wordSet[k] = v / float64(len(csvNames))
	}
	return wordSet, splitWords
}

func RemoveCommonWords(csvNames [][]string, threshold float64) {
	wordSet := make(map[string]int)

	jieba := jieba.NewJieba()
	defer jieba.Free()

	for _, row := range csvNames {
		words := jieba.Cut(row[1], true)
		tempWord := make(map[string]bool)
		for _, word := range words {
			tempWord[word] = true
		}
		for k, _ := range tempWord {
			if _, ok := wordSet[k]; ok {
				wordSet[k]++
			} else {
				wordSet[k] = 1
			}
		}
	}
	wordAboveThre := make([]string, 0)
	for k, v := range wordSet {
		if float64(v)/float64(len(csvNames)) >= threshold {
			wordAboveThre = append(wordAboveThre, k)
		}
	}
	for _, delWord := range wordAboveThre {
		for i := range csvNames {
			csvNames[i][1] = strings.Replace(csvNames[i][1], delWord, "", -1)
		}
	}
}

func BuildMapTree(rows [][]string, splitWords [][]string) (strMap map[string](uint64), iNodeMap map[uint64](*maptree.MapTreeNode), iNodeReverse map[uint64]string, iNodeToWords map[uint64]([]string)) {
	iNodeMap = make(map[uint64](*maptree.MapTreeNode))
	strMap = make(map[string](uint64))
	iNodeReverse = make(map[uint64]string)
	iNodeToWords = make(map[uint64]([]string))

	count := uint64(0)

	for rowI, row := range rows {
		if _, ok1 := strMap[row[1]]; ok1 {
			if _, ok2 := strMap[row[0]]; ok2 {
				tempNode := iNodeMap[strMap[row[0]]]
				tempNode = tempNode.ToRoot()
				tempName := iNodeReverse[tempNode.INode]
				if tempName != row[1] {
					fmt.Println("Error, simple name appear twice", row[0])
				}
			} else {
				parentNode := iNodeMap[strMap[row[1]]]
				parentNode = parentNode.ToRoot()
				newNode := maptree.NewMapTreeNodeWithParent(count, parentNode)
				count++
				strMap[row[0]] = newNode.INode
				iNodeReverse[newNode.INode] = row[0]
				iNodeMap[newNode.INode] = newNode
			}
		} else {
			if _, ok2 := strMap[row[0]]; ok2 {
				fmt.Println("Error, simple name appear twice", row[0])
			} else {
				strMap[row[1]] = count
				count++
				parentNode := maptree.NewMapTreeNode(strMap[row[1]])
				newNode := maptree.NewMapTreeNodeWithParent(count, parentNode)
				count++

				iNodeReverse[parentNode.INode] = row[1]
				iNodeReverse[newNode.INode] = row[0]
				strMap[row[0]] = newNode.INode
				iNodeMap[parentNode.INode] = parentNode
				iNodeMap[newNode.INode] = newNode
				iNodeToWords[parentNode.INode] = splitWords[rowI]
			}
		}
	}

	return
}

func DumpTree(strMap map[string](uint64), iNodeMap map[uint64](*maptree.MapTreeNode), iNodeToWords map[uint64]([]string), wordFreq map[string]float64, file_prefix string) {
	strMapFile, _ := os.Create(file_prefix + "str_map.pack")
	treeFile, _ := os.Create(file_prefix + "tree.pack")
	wordsFile, _ := os.Create(file_prefix + "words.pack")
	freqFile, _ := os.Create(file_prefix + "freq.pack")

	defer strMapFile.Close()
	defer treeFile.Close()
	defer wordsFile.Close()
	defer freqFile.Close()

	for k, v := range strMap {
		tempTreeNode := iNodeMap[v]
		strMapFile.WriteString(k + " " + strconv.FormatUint(v, 10) + "\n")
		strNode := strconv.FormatUint(tempTreeNode.INode, 10)
		if tempTreeNode.Parent != nil {
			strNode += " " + strconv.FormatUint(tempTreeNode.Parent.INode, 10)
		}
		treeFile.WriteString(strNode + "\n")
		if tempWords, ok := iNodeToWords[v]; ok {
			wordsFile.WriteString(strconv.FormatUint(v, 10) + " " + strings.Join(tempWords, " ") + "\n")
		}
	}
	for k, v := range wordFreq {
		freqFile.WriteString(k + " " + strconv.FormatFloat(v, 'f', -1, 64) + "\n")
	}
}

func LoadTree(file_prefix string) (strMap map[string](uint64), iNodeMap map[uint64](*maptree.MapTreeNode), iNodeReverse map[uint64]string, iNodeToWords map[uint64]([]string), wordFreq map[string]float64) {
	iNodeMap = make(map[uint64](*maptree.MapTreeNode))
	strMap = make(map[string](uint64))
	iNodeReverse = make(map[uint64]string)
	iNodeToWords = make(map[uint64]([]string))
	wordFreq = make(map[string]float64)

	strMapFile, _ := os.Open(file_prefix + "str_map.pack")
	defer strMapFile.Close()

	strFileBuf := bufio.NewReader(strMapFile)

	line, err := strFileBuf.ReadString('\n')

	for true {
		if len(line) > 0 {
			line = line[:len(line)-1]
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				iNode, _ := strconv.ParseUint(parts[1], 10, 64)
				strMap[parts[0]] = iNode
				iNodeReverse[iNode] = parts[0]
			}
		} else {
			break
		}
		if err == bufio.ErrFinalToken {
			break
		} else {
			line, err = strFileBuf.ReadString('\n')
		}
	}

	treeFile, _ := os.Open(file_prefix + "tree.pack")
	defer treeFile.Close()

	treeFileBuf := bufio.NewReader(treeFile)
	line, err = treeFileBuf.ReadString('\n')

	for true {
		if len(line) > 0 {
			line = line[:len(line)-1]
			parts := strings.Split(line, " ")
			parentNode := (*maptree.MapTreeNode)(nil)
			if len(parts) == 2 {
				iNode, _ := strconv.ParseUint(parts[1], 10, 64)
				if _, ok := iNodeMap[iNode]; ok {
					parentNode = iNodeMap[iNode]
				} else {
					parentNode = maptree.NewMapTreeNode(iNode)
					iNodeMap[iNode] = parentNode
				}
			}
			iNode, _ := strconv.ParseUint(parts[0], 10, 64)
			if _, ok := iNodeMap[iNode]; ok {
				iNodeMap[iNode].Parent = parentNode
			} else {
				newNode := maptree.NewMapTreeNodeWithParent(iNode, parentNode)
				iNodeMap[iNode] = newNode
			}
		} else {
			break
		}
		if err == bufio.ErrFinalToken {
			break
		} else {
			line, err = treeFileBuf.ReadString('\n')
		}
	}

	wordsFile, _ := os.Open(file_prefix + "words.pack")
	defer wordsFile.Close()

	wordsScanner := bufio.NewScanner(wordsFile)
	wordsScanner.Split(bufio.ScanLines)

	for wordsScanner.Scan() {
		line := wordsScanner.Text()
		if len(line) > 0 {
			words := strings.Split(line, " ")
			iNode, _ := strconv.ParseUint(words[0], 10, 64)
			iNodeToWords[iNode] = words[1:]
		}
	}

	freqFile, _ := os.Open(file_prefix + "freq.pack")
	defer freqFile.Close()

	freqScanner := bufio.NewScanner(freqFile)
	freqScanner.Split(bufio.ScanLines)

	for freqScanner.Scan() {
		line := freqScanner.Text()
		if len(line) > 0 {
			splits := strings.Split(line, " ")
			freq, _ := strconv.ParseFloat(splits[1], 64)
			wordFreq[splits[0]] = freq
		}
	}

	return
}
