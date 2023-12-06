package stddev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcDeviation(t *testing.T) {
	arr := []float32{
		0.1, 0.8, 0.1,
	}
	sd := CalcDeviation(arr)
	delta := 0.000001
	assert.InDelta(t, 0.4472136, sd, delta)

	arr2 := []float32{
		0, 1, 0,
	}
	sd = CalcDeviation(arr2)
	assert.InDelta(t, 0.0, sd, delta)

	arr = []float32{
		0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1,
	}
	sd = CalcDeviation(arr)
	assert.InDelta(t, 2.727636, sd, delta)

	arr = []float32{
		0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1,
	}
	sd = CalcDeviation(arr)
	assert.InDelta(t, 2.872281, sd, delta)

}
