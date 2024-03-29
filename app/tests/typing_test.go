package tests

import (
	"fmt"
	"go.uber.org/goleak"
	"reflect"
	"testing"
	"webApp/lib/eval"
	"webApp/lib/matrix"
)

func TestTyping(t *testing.T) {
	defer goleak.VerifyNone(t)

	if err := matrix.TypingMatrices(
		*GenerateMatrix(reflect.TypeOf(&eval.T1FS{}), reflect.TypeOf(eval.Interval{}), 200),
		*GenerateMatrix(reflect.TypeOf(&eval.T1FS{}), reflect.TypeOf(eval.Interval{}), 466)); err != nil {
		fmt.Println(err)
	}
}

func TestAggregateRatings(t *testing.T) {
	defer goleak.VerifyNone(t)

	if _, err := matrix.AggregateRatings([]matrix.Matrix{
		*GenerateMatrix(reflect.TypeOf(&eval.T1FS{}), reflect.TypeOf(eval.Interval{}), 200),
		*GenerateMatrix(reflect.TypeOf(&eval.T1FS{}), reflect.TypeOf(eval.Interval{}), 466)},
		[]eval.Evaluated{eval.Number(0.3), eval.Number(0.7)}); err != nil {
		fmt.Println(err)
	}
}
