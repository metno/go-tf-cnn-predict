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
}
