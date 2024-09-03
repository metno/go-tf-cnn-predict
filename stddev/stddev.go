package stddev

import (
	"math"
)

func CalcDeviation(x []float32) float64 {

	variance := math.Sqrt(math.Abs(CalcVariance(x)))

	return float64(variance)
}

/*
func CalcVariance(vector []float32) float64 {
	i := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}
	for j := range i {
		i[j] = i[j] / 8.0
	}
	x := vector
	var mean float32
	for j, val := range x {
		mean += i[j] * val
	}
	var variance float32
	for j, val := range x {
		variance += i[j] * i[j] * val
	}
	variance = variance - mean*mean
	return float64(variance)
}
*/

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
