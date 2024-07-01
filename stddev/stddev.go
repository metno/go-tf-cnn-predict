package stddev

import (
	"math"
)

func CalcDeviation(x []float32) float32 {

	variance := math.Sqrt(math.Abs(CalcVariance(x)))

	return float32(variance)
}

func CalcVariance(arr []float32) float64 {
	var mean float32 = 0
	for _, value := range arr {
		mean += value
	}
	mean /= float32(len(arr))

	var variance float32 = 0
	for _, value := range arr {
		difference := value - mean
		variance += difference * difference
	}
	variance /= float32(len(arr))

	return float64(variance)
}
