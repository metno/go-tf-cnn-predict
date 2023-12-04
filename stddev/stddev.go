package stddev

import (
	"math"
)

func argmax(arr []float32) int {
	var big float32 = -1.0
	var idx int = -1

	for i := 0; i < len(arr); i++ {
		if arr[i] >= big {
			big = arr[i]
			idx = i
		}
	}
	return idx
}
func argmin(arr []float32) int {
	var small float32 = 1.0
	var idx int = -1

	for i := 0; i < len(arr); i++ {
		if arr[i] <= small {
			small = arr[i]
			idx = i
		}
	}
	return idx
}

func CalcDeviation(num []float32) float32 {
	i := make([]float32, len(num))
	var mean float32 = 0
	for j := 0; j < len(num); j++ {
		i[j] = float32(j)
	}
	for j := 0; j < len(num); j++ {
		mean += num[j] * i[j]
	}

	var variance float32 = 0.0
	for j := 0; j < len(num); j++ {
		variance += num[j] * i[j] * i[j]
	}
	variance -= mean * mean

	// Take math.Abs, sqrt(-0) = Nan ..
	variance = float32(math.Sqrt(math.Abs(float64(variance))))

	return variance
}
