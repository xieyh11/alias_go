package strsim

const (
	MapSimiliarityLevelOne   = iota // province level
	MapSimiliarityLevelTwo          // city level
	MapSimiliarityLevelThree        // district level
)

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
	if segRes[0] == 0.0 && level == MapSimiliarityLevelOne {
		return 0
	}
	if segRes[1] == 0.0 && (level == MapSimiliarityLevelOne || level == MapSimiliarityLevelTwo) {
		return 0
	}
	if segRes[2] == 0.0 {
		return 0
	}
	res := 1.0
	for i := 3; i < len(segRes); i++ {
		res *= segRes[i]
	}
	return res
}
