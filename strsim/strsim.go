package strsim

import (
	"../floatvector"
	"github.com/cpmech/gosl/graph"
)

func MaxCommonScore(str1, str2 string, addDelScore, matchScore, replaceScore float64) float64 {
	return RunesMaxCommonScore([]rune(str1), []rune(str2), addDelScore, matchScore, replaceScore)
}

func RunesMaxCommonScore(str1, str2 []rune, addDelScore, matchScore, replaceScore float64) float64 {
	if len(str1) == 0 {
		return addDelScore * float64(len(str2))
	}
	if len(str2) == 0 {
		return addDelScore * float64(len(str1))
	}

	lenstr1, lenstr2 := len(str1), len(str2)
	res := make([][]float64, lenstr1)
	for i := range res {
		res[i] = make([]float64, lenstr2)
	}
	if str1[lenstr1-1] == str2[lenstr2-1] {
		res[lenstr1-1][lenstr2-1] = matchScore
	} else {
		res[lenstr1-1][lenstr2-1] = replaceScore
	}

	//last row
	for j := lenstr2 - 2; j >= 0; j-- {
		initScore := matchScore
		if str1[lenstr1-1] != str2[j] {
			initScore = replaceScore
		}
		initScore += addDelScore * float64(lenstr2-1-j)
		if initScore > (res[lenstr1-1][j+1] + addDelScore) {
			res[lenstr1-1][j] = initScore
		} else {
			res[lenstr1-1][j] = res[lenstr1-1][j+1] + addDelScore
		}
	}

	//last column
	for i := lenstr1 - 2; i >= 0; i-- {
		initScore := matchScore
		if str1[i] != str2[lenstr2-1] {
			initScore = replaceScore
		}
		initScore += addDelScore * float64(lenstr1-1-i)
		if initScore > (res[i+1][lenstr2-1] + addDelScore) {
			res[i][lenstr2-1] = initScore
		} else {
			res[i][lenstr2-1] = res[i+1][lenstr2-1] + addDelScore
		}
	}

	for row, col := lenstr1-2, lenstr2-2; row >= 0 && col >= 0; row, col = row-1, col-1 {
		initScore := matchScore
		if str1[row] != str2[col] {
			initScore = replaceScore
		}
		initScore += res[row+1][col+1]
		if initScore < (res[row+1][col] + addDelScore) {
			initScore = res[row+1][col] + addDelScore
		}
		if initScore < (res[row][col+1] + addDelScore) {
			initScore = res[row][col+1] + addDelScore
		}

		res[row][col] = initScore

		//row
		for j := col - 1; j >= 0; j-- {
			initScore := matchScore
			if str1[row] != str2[j] {
				initScore = replaceScore
			}
			initScore += res[row+1][j+1]
			if initScore < (res[row][j+1] + addDelScore) {
				initScore = res[row][j+1] + addDelScore
			}
			if initScore < (res[row+1][j] + addDelScore) {
				initScore = res[row+1][j] + addDelScore
			}
			res[row][j] = initScore
		}

		//column
		for i := row - 1; i >= 0; i-- {
			initScore := matchScore
			if str1[i] != str2[col] {
				initScore = replaceScore
			}
			initScore += res[i+1][col+1]
			if initScore < (res[i+1][col] + addDelScore) {
				initScore = res[i+1][col] + addDelScore
			}
			if initScore < (res[i][col+1] + addDelScore) {
				initScore = res[i][col+1] + addDelScore
			}
			res[i][col] = initScore
		}
	}
	minS := lenstr1
	maxS := lenstr2
	if minS > maxS {
		minS, maxS = maxS, minS
	}
	maxPossible := float64(minS)*matchScore + float64(maxS-minS)*addDelScore
	return res[0][0] / maxPossible
}

func SplitCommonScore(str1, base []string, wordWeight map[string]float64) float64 {
	var mnk graph.Munkres
	rows, cols := len(str1), len(base)
	costMatrix := make([][]float64, rows)
	for i := range costMatrix {
		costMatrix[i] = make([]float64, cols)
		for j := range costMatrix[i] {
			costMatrix[i][j] = -RunesMaxCommonScore([]rune(str1[i]), []rune(base[j]), 0, 1, 0)
		}
	}
	mnk.Init(rows, cols)
	mnk.SetCostMatrix(costMatrix)
	mnk.Run()
	totalFreq := float64(0)
	currentScore := float64(0)
	for i := 0; i < rows; i++ {
		j := mnk.Links[i]
		if j != -1 {
			totalFreq += wordWeight[base[j]]
			currentScore += wordWeight[base[j]] * (-costMatrix[i][j])
		}
	}
	return currentScore / totalFreq
}

func SplitVectorDis(str1, base []string, str1Vec [][]float64, wordWeight map[string]float64, wordVector map[string][]float64) float64 {
	var mnk graph.Munkres
	rows, cols := len(str1), len(base)
	costMatrix := make([][]float64, rows)
	baseVec := make([][]float64, 0)
	for i := range base {
		baseVec = append(baseVec, wordVector[base[i]])
	}
	for i := range costMatrix {
		costMatrix[i] = make([]float64, cols)
		for j := range costMatrix[i] {
			costMatrix[i][j] = floatvector.VectorsDis(str1Vec[i], baseVec[j])
		}
	}
	mnk.Init(rows, cols)
	mnk.SetCostMatrix(costMatrix)
	mnk.Run()
	totalFreq := float64(0)
	currentScore := float64(0)
	for i := 0; i < rows; i++ {
		j := mnk.Links[i]
		if j != -1 {
			totalFreq += wordWeight[base[j]]
			currentScore += wordWeight[base[j]] * (-costMatrix[i][j])
		}
	}
	return currentScore / totalFreq
}
