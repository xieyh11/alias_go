package searchtree

import (
	"container/heap"
	"strings"

	"../strsim"
)

func subSearch(names []string, scores []float64, name string, count chan int) {
	for strI := range names {
		scores[strI] = strsim.RunesMaxCommonScore([]rune(names[strI]), []rune(name), 0, 1, 0)
	}
	count <- 1
}

func SearchStrings(names []string, name string, threads int) []float64 {
	scores := make([]float64, len(names))
	per_size := len(names) / threads
	count := make(chan int)
	for i := 0; i < threads; i++ {
		if i == threads-1 {
			go subSearch(names[i*per_size:], scores[i*per_size:], name, count)
		} else {
			go subSearch(names[i*per_size:i*per_size+per_size], scores[i*per_size:i*per_size+per_size], name, count)
		}
	}
	defer close(count)
	returnThread := 0
	for {
		if _, ok := <-count; ok {
			returnThread++
		}
		if returnThread == threads {
			break
		}
	}
	return scores
}

func subSearchSplitStrings(names []string, strMap map[string]uint64, iNodeToWords map[uint64]([]string), wordFreq map[string]float64, name []string, scores []float64, count chan int) {
	for strI := range names {
		if v, ok := iNodeToWords[strMap[names[strI]]]; ok {
			scores[strI] = strsim.SplitCommonScore(name, v, wordFreq)
		} else {
			scores[strI] = strsim.RunesMaxCommonScore([]rune(names[strI]), []rune(strings.Join(name, "")), 0, 1, 0)
		}
	}
	count <- 1
}

func SearchSplitStrings(names []string, strMap map[string]uint64, iNodeToWords map[uint64]([]string), wordFreq map[string]float64, name []string, threads int) []float64 {
	scores := make([]float64, len(names))
	per_size := len(names) / threads
	count := make(chan int)
	for i := 0; i < threads; i++ {
		if i == threads-1 {
			go subSearchSplitStrings(names[i*per_size:], strMap, iNodeToWords, wordFreq, name, scores[i*per_size:], count)
		} else {
			go subSearchSplitStrings(names[i*per_size:i*per_size+per_size], strMap, iNodeToWords, wordFreq, name, scores[i*per_size:i*per_size+per_size], count)
		}
	}
	defer close(count)
	returnThread := 0
	for {
		if _, ok := <-count; ok {
			returnThread++
		}
		if returnThread == threads {
			break
		}
	}
	return scores
}

type ScoreIndex struct {
	Idx   int
	Score float64
}

func NewScoreIndex(idx int, score float64) *ScoreIndex {
	res := new(ScoreIndex)
	res.Idx = idx
	res.Score = score
	return res
}

type priorityQueue []*ScoreIndex

func (p priorityQueue) Len() int { return len(p) }
func (p priorityQueue) Less(i, j int) bool {
	return p[i].Score < p[j].Score
}
func (p priorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *priorityQueue) Push(x interface{}) {
	item := x.(*ScoreIndex)
	*p = append(*p, item)
}

func (p *priorityQueue) Pop() interface{} {
	old := *p
	item := old[len(old)-1]
	*p = old[0 : len(old)-1]
	return item
}

func TopKScores(scores []float64, k int) []int {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)

	for i := range scores {
		if len(pq) < k {
			heap.Push(&pq, NewScoreIndex(i, scores[i]))
		} else {
			if pq[0].Score < scores[i] {
				heap.Pop(&pq)
				heap.Push(&pq, NewScoreIndex(i, scores[i]))
			}
		}
	}
	res := make([]int, 0)
	for len(pq) > 0 {
		item := heap.Pop(&pq).(*ScoreIndex)
		res = append([]int{item.Idx}, res...)
	}
	return res
}
