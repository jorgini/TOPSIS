package tests

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func TestEmptyValues(t *testing.T) {
	defer goleak.VerifyNone(t)

	m := matrix.NewMatrix(3, 3)
	err := m.SetRatings([][]eval.Evaluated{
		{eval.Number(0), eval.Number(0).ConvertToAIFS(v.Triangle), eval.Number(0).ConvertToInterval()},
		{eval.Number(0).ConvertToT1FS(v.Trapezoid), eval.Number(0).ConvertToInterval(), eval.Number(0)},
		{eval.Number(0).ConvertToAIFS(v.Trapezoid), eval.Number(0), eval.Number(0).ConvertToInterval()},
	})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.SetCriteria([]matrix.Criterion{{Weight: eval.Rating{eval.Number(5.5)}, TypeOfCriteria: v.Benefit},
		{Weight: eval.Rating{eval.Number(3.2)}, TypeOfCriteria: v.Cost},
		{Weight: eval.Rating{eval.Number(8)}, TypeOfCriteria: v.Benefit}})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = matrix.TypingMatrices(*m)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.Normalization(v.NormalizeWithSum, v.NormalizeWithSum)

	assert.Equal(t, v.EmptyValues, err)
}

func TestEmptyWeights(t *testing.T) {
	defer goleak.VerifyNone(t)

	m := matrix.NewMatrix(3, 3)
	err := m.SetRatings([][]eval.Evaluated{
		{eval.Number(4), eval.Number(0).ConvertToAIFS(v.Triangle), eval.Number(8.4).ConvertToInterval()},
		{eval.Number(3).ConvertToT1FS(v.Trapezoid), eval.Number(6.1).ConvertToInterval(), eval.Number(7)},
		{eval.Number(5.5).ConvertToAIFS(v.Trapezoid), eval.Number(4), eval.Number(0.4).ConvertToInterval()},
	})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.SetCriteria([]matrix.Criterion{{Weight: eval.Rating{eval.Number(0)}, TypeOfCriteria: v.Benefit},
		{Weight: eval.Rating{eval.Number(0)}, TypeOfCriteria: v.Cost},
		{Weight: eval.Rating{eval.Number(0)}, TypeOfCriteria: v.Benefit}})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = matrix.TypingMatrices(*m)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.Normalization(v.NormalizeWithSum, v.NormalizeWithSum)

	assert.Equal(t, v.EmptyValues, err)
}

func TestWrongSizeValue(t *testing.T) {
	defer goleak.VerifyNone(t)

	m := matrix.NewMatrix(3, 3)
	err := m.SetRatings([][]eval.Evaluated{
		{eval.Number(4), eval.NumbersMin, eval.Number(0), eval.Number(8.4)},
		{eval.Number(3), eval.Number(6.1), eval.Number(7)},
		{eval.Number(5.5), eval.Number(4), eval.Number(0.4)},
	})

	assert.Equal(t, v.InvalidSize, err)
}

func TestWrongSizeWeights(t *testing.T) {
	defer goleak.VerifyNone(t)

	m := matrix.NewMatrix(3, 3)
	err := m.SetRatings([][]eval.Evaluated{
		{eval.Number(4), eval.Number(0), eval.Number(8.4)},
		{eval.Number(3), eval.Number(6.1), eval.Number(7)},
		{eval.Number(5.5), eval.Number(4), eval.Number(0.4)},
	})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.SetCriteria([]matrix.Criterion{{Weight: eval.Rating{eval.Number(5.5)}, TypeOfCriteria: v.Benefit},
		{Weight: eval.Rating{eval.Number(3.2)}, TypeOfCriteria: v.Cost},
		{Weight: eval.Rating{eval.Number(8)}, TypeOfCriteria: v.Benefit},
		{Weight: eval.Rating{eval.Number(7)}, TypeOfCriteria: v.Benefit}})

	assert.Equal(t, v.InvalidSize, err)
}
