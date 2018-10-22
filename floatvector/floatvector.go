package floatvector

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func AddVectorsInPlace(dst, src []float64) {
	if len(dst) != len(src) {
		panic(fmt.Sprintln("Unequal length in Vector Sum: ", len(dst), len(src)))
	}
	for i := range src {
		dst[i] += src[i]
	}
}

func VectorsDis(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		panic(fmt.Sprintln("Unequal length in Vector Dis: ", len(v1), len(v2)))
	}
	sum := 0.0
	for i := range v1 {
		sum += (v1[i] - v2[i]) * (v1[i] - v2[i])
	}
	return math.Sqrt(sum)
}

func VectorToString(v []float64) string {
	res := make([]string, len(v))
	for i := range v {
		res[i] = strconv.FormatFloat(v[i], 'f', -1, 64)
	}
	return strings.Join(res, " ")
}
