package matrix

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
	"webApp/lib/eval"
	v "webApp/lib/variables"
)

func TestEmptyValues(t *testing.T) {
	defer goleak.VerifyNone(t)

	m := NewMatrix(3, 3)
	err := m.SetRatings([][]eval.Evaluated{
		{eval.Number(0), eval.Number(0).ConvertToAIFS(v.Triangle), eval.Number(0).ConvertToInterval()},
		{eval.Number(0).ConvertToT1FS(v.Trapezoid), eval.Number(0).ConvertToInterval(), eval.Number(0)},
		{eval.Number(0).ConvertToAIFS(v.Trapezoid), eval.Number(0), eval.Number(0).ConvertToInterval()},
	})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.SetCriteria([]Criterion{{Weight: eval.Rating{eval.Number(5.5)}, TypeOfCriteria: v.Benefit},
		{Weight: eval.Rating{eval.Number(3.2)}, TypeOfCriteria: v.Cost},
		{Weight: eval.Rating{eval.Number(8)}, TypeOfCriteria: v.Benefit}})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = TypingMatrices(*m)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.Normalization(v.NormalizeWithSum, v.NormalizeWithSum)

	assert.Equal(t, v.EmptyValues, err)
}

func TestEmptyWeights(t *testing.T) {
	defer goleak.VerifyNone(t)

	m := NewMatrix(3, 3)
	err := m.SetRatings([][]eval.Evaluated{
		{eval.Number(4), eval.Number(0).ConvertToAIFS(v.Triangle), eval.Number(8.4).ConvertToInterval()},
		{eval.Number(3).ConvertToT1FS(v.Trapezoid), eval.Number(6.1).ConvertToInterval(), eval.Number(7)},
		{eval.Number(5.5).ConvertToAIFS(v.Trapezoid), eval.Number(4), eval.Number(0.4).ConvertToInterval()},
	})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.SetCriteria([]Criterion{{Weight: eval.Rating{eval.Number(0)}, TypeOfCriteria: v.Benefit},
		{Weight: eval.Rating{eval.Number(0)}, TypeOfCriteria: v.Cost},
		{Weight: eval.Rating{eval.Number(0)}, TypeOfCriteria: v.Benefit}})

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = TypingMatrices(*m)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = m.Normalization(v.NormalizeWithSum, v.NormalizeWithSum)

	assert.Equal(t, v.EmptyValues, err)
}

//func TestMatrix_Normalization(t *testing.T) {
//	testCases := []struct {
//		name          string
//		input         Matrix
//		optNormValue  v.Variants
//		optNormWeight v.Variants
//		expected      Matrix
//	}{
//		{
//			name: "Numbers",
//			input: Matrix{
//				Data: []Alternative{
//					{
//						Grade:           []eval.Rating{{eval.Number(5.5)}, {eval.Number(7)}, {eval.Number(3.7)}},
//						CountOfCriteria: 3,
//					},
//					{
//						Grade:           []eval.Rating{{eval.Number(4.8)}, {eval.Number(8.2)}, {eval.Number(6.6)}},
//						CountOfCriteria: 3,
//					},
//					{
//						Grade:           []eval.Rating{{eval.Number(9.1)}, {eval.Number(3)}, {eval.Number(2.8)}},
//						CountOfCriteria: 3,
//					},
//				},
//				Criteria: []Criterion{
//					{
//						Weight:         eval.Rating{eval.Number(5)},
//						TypeOfCriteria: v.Benefit,
//					},
//					{
//						Weight:         eval.Rating{eval.Number(1)},
//						TypeOfCriteria: v.Cost,
//					},
//					{
//						Weight:         eval.Rating{eval.Number(7)},
//						TypeOfCriteria: v.Cost,
//					},
//				},
//				CountAlternatives: 3,
//				CountCriteria:     3,
//				HighType:          "eval.Number",
//				FormFs:            v.None,
//			},
//			optNormValue:  v.NormalizeWithSum,
//			optNormWeight: v.NormalizeWithSum,
//			expected: Matrix{
//				Data: []Alternative{
//					{
//						Grade:           []eval.Rating{{eval.Number(0.471)}, {eval.Number(0.625)}, {eval.Number(0.459)}},
//						CountOfCriteria: 3,
//					},
//					{
//						Grade:           []eval.Rating{{eval.Number(0.411)}, {eval.Number(0.733)}, {eval.Number(0.818)}},
//						CountOfCriteria: 3,
//					},
//					{
//						Grade:           []eval.Rating{{eval.Number(0.780)}, {eval.Number(0.268)}, {eval.Number(0.347)}},
//						CountOfCriteria: 3,
//					},
//				},
//				Criteria: []Criterion{
//					{
//						Weight: eval.Rating{eval.Number(0.385)},
//						TypeOfCriteria: v.Benefit,
//					},
//					{
//						Weight: eval.Rating{eval.Number(0.077)},
//						TypeOfCriteria: v.Cost,
//					},{
//						Weight: eval.Rating{eval.Number(0.538)},
//						TypeOfCriteria: v.Cost,
//					},
//				},
//				CountCriteria: 3,
//				CountAlternatives: 3,
//				HighType: "eval.Number",
//				FormFs: v.None,
//			},
//		},
//	}
//
//	for _, tt := range testCases {
//		t.Run(tt.name, func(t *testing.T) {
//			defer goleak.VerifyNone(t)
//			if err := tt.input.Normalization(tt.optNormValue, tt.optNormWeight); err != nil {
//				t.Errorf(err.Error())
//			}
//
//			for
//		})
//	}
//}
