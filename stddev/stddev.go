package stddev

import (
	"math"
)

func CalcDeviation(x []float32) float32 {
	//fmt.Printf("x: %v\n", x)
	i := make([]float32, len(x))
	var mean float32 = 0
	for j := 0; j < len(x); j++ {
		i[j] = float32(j) / float32((len(x) - 1.0))
	}
	var variance float32 = 0.0
	//fmt.Printf("i: %v\n", i)
	for j := 0; j < len(x); j++ {
		mean += x[j] * i[j]
		variance = variance + (x[j] * i[j] * i[j])
	}
	variance -= (mean * mean)

	// Take math.Abs, sqrt(-0) = Nan ..
	variance = float32((len(x) - 1.0)) * float32(math.Sqrt(math.Abs(float64(variance))))

	return variance
}
