package stddev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var delta = 0.001

func TestCalcDeviation(t *testing.T) {
	arr := []float32{
		0.0, 0.1, 0.0,
	}
	sd := CalcDeviation(arr)
	assert.InDelta(t, 0.0471, sd, delta)

	arr = []float32{
		1.0, 1.0, 1.0,
	}

	sd = CalcDeviation(arr)
	assert.InDelta(t, 0.0, sd, delta)
}

func TestCalcVariance(t *testing.T) {
	arr := []float32{
		0.0, 1.0, 0.0,
	}
	variance := CalcVariance(arr)
	expectedVariance := 0.22222 // since all elements of array are equal, variance should be zero
	assert.InDelta(t, expectedVariance, variance, delta)

	arr = []float32{
		1.0, 1.0, 1.0,
	}
	variance = CalcVariance(arr)
	expectedVariance = 0.0 // since all elements of array are equal, variance should be zero
	assert.InDelta(t, expectedVariance, variance, delta)
}
