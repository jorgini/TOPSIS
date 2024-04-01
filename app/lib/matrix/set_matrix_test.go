package matrix

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
	"webApp/lib/eval"
	v "webApp/lib/variables"
)

func TestWrongSizeValue(t *testing.T) {
	defer goleak.VerifyNone(t)

	m := NewMatrix(3, 3)
	err := m.setRatings([][]eval.Evaluated{
		{eval.Number(4), eval.NumbersMin, eval.Number(0), eval.Number(8.4)},
		{eval.Number(3), eval.Number(6.1), eval.Number(7)},
		{eval.Number(5.5), eval.Number(4), eval.Number(0.4)},
	})

	assert.Equal(t, v.InvalidSize, err)
}

func TestWrongSizeWeights(t *testing.T) {
	defer goleak.VerifyNone(t)

	m := NewMatrix(3, 3)
	err := m.setRatings([][]eval.Evaluated{
		{eval.Number(4), eval.Number(0), eval.Number(8.4)},
		{eval.Number(3), eval.Number(6.1), eval.Number(7)},
		{eval.Number(5.5), eval.Number(4), eval.Number(0.4)},
	})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.SetCriteria([]Criterion{{Weight: eval.Rating{eval.Number(5.5)}, TypeOfCriteria: v.Benefit},
		{Weight: eval.Rating{eval.Number(3.2)}, TypeOfCriteria: v.Cost},
		{Weight: eval.Rating{eval.Number(8)}, TypeOfCriteria: v.Benefit},
		{Weight: eval.Rating{eval.Number(7)}, TypeOfCriteria: v.Benefit}})

	assert.Equal(t, v.InvalidSize, err)
}
