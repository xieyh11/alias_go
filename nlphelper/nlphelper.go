package nlphelper

type WordsSegment interface {
	Segment(str string) []string
}

type Word2Vector interface {
	ToVector(str string) []float64
}

const (
	NlpUsingJieba = iota
	NlpUsingHanlp
	NlpUsingGaode
)
