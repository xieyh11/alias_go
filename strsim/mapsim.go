package strsim

import (
	"regexp"
	"strings"
)

const (
	MapSimiliarityLevelOne   = iota // province level
	MapSimiliarityLevelTwo          // city level
	MapSimiliarityLevelThree        // district level
)

var pathKeywords = []string{"片区", "大道", "路", "街", "道", "村", "里"}
var numKeywords = []string{"单元", "号楼", "弄", "号", "栋", "楼", "幢"}

func isPath(rowStr string, idx, size int) bool {
	fullPath := false
	for _, k := range pathKeywords {
		if idx+size-len(k) >= 0 {
			if k == rowStr[idx+size-len(k):idx+size] {
				fullPath = true
			}
		}
	}
	if fullPath {
		return true
	}
	appendPath := false
	if idx+size < len(rowStr) {
		for _, k := range pathKeywords {
			endIdx := idx + size + len(k)
			if endIdx > len(rowStr) {
				endIdx = len(rowStr)
			}
			if k == rowStr[idx+size:endIdx] {
				appendPath = true
			}
		}
	}
	return appendPath
}

var provinceSuffix = []string{"维吾尔自治区", "回族自治区", "特别行政区", "壮族自治区", "自治区", "直辖市", "省", "市"}
var citySuffix = []string{"市", "区"}
var adSuffix = []string{"市", "县", "区"}

func removeAdDiv(rawStr, name string, suffix []string) string {
	if name != "" {
		ni := strings.Index(rawStr, name)
		if ni == -1 {
			hasChange := false
			for _, suffixI := range suffix {
				if strings.HasSuffix(name, suffixI) {
					name = name[:len(name)-len(suffixI)]
					hasChange = true
					break
				}
			}
			if hasChange {
				ni = strings.Index(rawStr, name)
				if ni != -1 {
					if !isPath(rawStr, ni, len(name)) {
						rawStr = strings.Replace(rawStr, name, "", 1)
					}
				}
			}
		} else {
			rawStr = strings.Replace(rawStr, name, "", 1)
		}
	}
	return rawStr
}

func findPath(rowStr string) (int, int) {
	for _, k := range pathKeywords {
		idx := strings.Index(rowStr, k)
		if idx != -1 {
			return idx, len(k)
		}
	}
	return -1, 0
}

func removeStrByIndex(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	return s[:start] + s[end:]
}

func getSubString(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}

func MapExtractInfo(rawStr string, search []string) []string {
	res := make([]string, 0)
	pname := search[0]
	res = append(res, pname)
	rawStr = removeAdDiv(rawStr, pname, provinceSuffix)

	cityname := search[1]
	res = append(res, cityname)
	rawStr = removeAdDiv(rawStr, cityname, citySuffix)

	adname := search[2]
	res = append(res, adname)
	rawStr = removeAdDiv(rawStr, adname, adSuffix)

	pathI, pathL := findPath(rawStr)
	if pathI != -1 {
		res = append(res, rawStr[:pathI+pathL])
		rawStr = rawStr[pathI+pathL:]
		pathI, pathL = findPath(rawStr)
		if pathI != -1 {
			res[len(res)-1] += rawStr[:pathI+pathL]
			rawStr = rawStr[pathI+pathL:]
		}
	} else {
		pathI, pathL := findPath(search[3])
		if pathI != -1 {
			res = append(res, search[3][:pathI+pathL])
		} else {
			res = append(res, "")
		}
	}

	numReg, _ := regexp.Compile(`[a-zA-Z]{0,1}\d+`)
	numIdxs := numReg.FindIndex([]byte(rawStr))
	if numIdxs != nil {
		res = append(res, rawStr[numIdxs[0]:numIdxs[1]])
		hasKey := false
		for _, k := range numKeywords {
			endIdx := numIdxs[1] + len(k)
			if getSubString(rawStr, numIdxs[1], endIdx) == k {
				rawStr = removeStrByIndex(rawStr, numIdxs[0], endIdx)
				hasKey = true
				break
			}
		}
		if !hasKey {
			rawStr = removeStrByIndex(rawStr, numIdxs[0], numIdxs[1])
		}
	} else {
		numIdxs := numReg.FindIndex([]byte(search[3]))
		if numIdxs != nil {
			res = append(res, search[3][numIdxs[0]:numIdxs[1]])
		} else {
			res = append(res, "")
		}
	}
	numIdxs = numReg.FindIndex([]byte(rawStr))
	numStrs := make([]string, 0)
	for numIdxs != nil {
		numStrs = append(numStrs, rawStr[numIdxs[0]:numIdxs[1]])
		hasKey := false
		for _, k := range numKeywords {
			endIdx := numIdxs[1] + len(k)
			if getSubString(rawStr, numIdxs[1], endIdx) == k {
				rawStr = removeStrByIndex(rawStr, numIdxs[0], endIdx)
				hasKey = true
				break
			}
		}
		if !hasKey {
			rawStr = removeStrByIndex(rawStr, numIdxs[0], numIdxs[1])
		}
		numIdxs = numReg.FindIndex([]byte(rawStr))
	}
	res = append(res, strings.Join(numStrs, "-"))

	res = append(res, rawStr)

	return res
}

var mapWeights = []float64{0.8, 0.1, 0.5, 0.5}

func tripleMax(a, b, c int) int {
	res := a
	if res < b {
		res = b
	}
	if res < c {
		res = c
	}
	return res
}

func MaxSubSeq(seq1, seq2 []string) int {
	rows, cols := len(seq1), len(seq2)
	if rows == 0 || cols == 0 {
		return 0
	}
	res := make([][]int, rows+1)
	for i := range res {
		res[i] = make([]int, cols+1)
	}
	for i, j := rows-1, cols-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		ij := 0
		if seq1[i] == seq2[j] {
			ij = 1
		}
		res[i][j] = tripleMax(ij+res[i+1][j+1], res[i][j+1], res[i+1][j])
		for c := j - 1; c >= 0; c-- {
			ic := 0
			if seq1[i] == seq2[c] {
				ic = 1
			}
			res[i][c] = tripleMax(ic+res[i+1][c+1], res[i][c+1], res[i+1][c])
		}
		for r := i - 1; r >= 0; r-- {
			rj := 0
			if seq1[r] == seq2[j] {
				rj = 1
			}
			res[r][j] = tripleMax(rj+res[r+1][j+1], res[r][j+1], res[r+1][j])
		}
	}
	return res[0][0]
}

func MapSegmentSimiliarity(map1, map2 []string, level int) float64 {
	if len(map1) != len(map2) {
		return 0
	}
	segL := len(map1)
	segRes := make([]float64, segL)
	for i := range segRes {
		switch {
		case i < 3:
			fallthrough
		case i == 4:
			if len(map1[i]) == 0 || len(map2[i]) == 0 {
				segRes[i] = 1
			} else {
				if map1[i] == map2[i] {
					segRes[i] = 1
				} else {
					segRes[i] = 0
				}
			}
		case i == 3:
			segRes[i] = RunesMaxCommonScore([]rune(map1[i]), []rune(map2[i]), 0, 1, 0)
		case i == 5:
			if len(map1[i]) == 0 || len(map2[i]) == 0 {
				segRes[i] = 1
			} else {
				nums1 := strings.Split(map1[i], "-")
				nums2 := strings.Split(map2[i], "-")
				maxL := MaxSubSeq(nums1, nums2)
				segRes[i] = float64(maxL) / float64(len(nums1)+len(nums2)-maxL)
			}
		default:
			segRes[i] = RunesMaxCommonScore([]rune(map1[i]), []rune(map2[i]), 0, 1, 0)
		}
	}
	if segRes[0] == 0.0 {
		return 0
	}
	if segRes[1] == 0.0 && (level == MapSimiliarityLevelTwo || level == MapSimiliarityLevelThree) {
		return 0
	}
	if segRes[2] == 0.0 && level == MapSimiliarityLevelThree {
		return 0
	}
	res := 0.0
	weights := 0.0

	for i := 3; i < len(segRes); i++ {
		if map1[i] != "" && map2[i] != "" {
			res += segRes[i] * mapWeights[i-3]
			weights += mapWeights[i-3]
		}
	}
	return res / weights
}
