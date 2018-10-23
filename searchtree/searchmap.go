package searchtree

import (
	"../buildtree"
	"../strsim"
)

func subSearchMap(names []string, mapTree buildtree.MapTree, strMessage buildtree.StrMessage, name []string, scores []float64, level int, count chan int) {
	for strI := range names {
		scores[strI] = strsim.MapSegmentSimiliarity(strMessage.INodeToWords[mapTree.StrToINode[names[strI]]], name, level)
	}
	count <- 1
}

func SearchMap(names []string, mapTree buildtree.MapTree, strMessage buildtree.StrMessage, name []string, level int, threads int) []float64 {
	scores := make([]float64, len(names))
	per_size := len(names) / threads
	count := make(chan int)
	for i := 0; i < threads; i++ {
		if i == threads-1 {
			go subSearchMap(names[i*per_size:], mapTree, strMessage, name, scores[i*per_size:], level, count)
		} else {
			go subSearchMap(names[i*per_size:i*per_size+per_size], mapTree, strMessage, name, scores[i*per_size:i*per_size+per_size], level, count)
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
