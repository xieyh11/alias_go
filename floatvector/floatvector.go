package floatvector

func AddVectorsInPlace(dst, src []float64) {
	for i := range src {
		dst[i] += src[i]
	}
}
