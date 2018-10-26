package strsim

import (
	"regexp"
	"strings"

	"../stringhelper"
)

const (
	MapSimiliarityLevelOne   = iota // province level
	MapSimiliarityLevelTwo          // city level
	MapSimiliarityLevelThree        // district level
)

var pathKeywords = []string{"片区", "大道", "路", "街", "道", "村"}
var numKeywords = []string{"单元", "号楼", "弄", "号", "栋", "楼"}

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
						strings.Replace(rawStr, name, "", 1)
					}
				}
			}
		} else {
			strings.Replace(rawStr, name, "", 1)
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
	}

	numReg, _ := regexp.Compile(`[a-zA-Z]{0,1}\d+`)
	numIdxs := numReg.FindAllIndex([]byte(rawStr), -1)
	if numIdxs != nil {
		res = append(res, rawStr[numIdxs[0][0]:numIdxs[0][1]])
	} else {
		res = append(res, "")
	}
	if len(numIdxs) > 1 {
		numStrs := make([]string, 0)
		for i := 1; i < len(numIdxs); i++ {
			numStrs = append(numStrs, rawStr[numIdxs[i][0]:numIdxs[i][1]])
		}
		res = append(res, strings.Join(numStrs, "-"))
	} else {
		res = append(res, "")
	}

	for _, idx := range numIdxs {
		if idx[1] < len(rawStr) {
			for _, k := range numKeywords {
				endIdx := idx[1] + len(k)
				if endIdx > len(rawStr) {
					endIdx = len(rawStr)
				}
				if rawStr[idx[1]:endIdx] == k {
					rawStr = rawStr[:idx[0]] + rawStr[endIdx:]
					break
				}
			}
		}
	}

	rawStr = stringhelper.RemovePunc(rawStr)
	res = append(res, rawStr)

	return res
}

var mapWeights = []float64{0.8, 0.1, 0.5, 0.5}

func MapSegmentSimiliarity(map1, map2 []string, level int) float64 {
	if len(map1) != len(map2) {
		return 0
	}
	segL := len(map1)
	segRes := make([]float64, segL)
	for i := range segRes {
		if i < 3 {
			if len(map1[i]) == 0 || len(map2[i]) == 0 {
				segRes[i] = 1
			} else {
				if map1[i] == map2[i] {
					segRes[i] = 1
				} else {
					segRes[i] = 0
				}
			}
		} else {
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
