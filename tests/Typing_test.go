package tests

import (
	"fmt"
	topsis "github.com/jorka/TOPSIS/lib"
	"go.uber.org/goleak"
	"reflect"
	"testing"
)

func TestTyping(t *testing.T) {
	defer goleak.VerifyNone(t)

	if err := topsis.TypingMatrices(
		GenerateMatrix(reflect.TypeOf(&topsis.T1FS{}), reflect.TypeOf(topsis.Interval{}), 200),
		GenerateMatrix(reflect.TypeOf(&topsis.T1FS{}), reflect.TypeOf(topsis.Interval{}), 466)); err != nil {
		fmt.Println(err)
	}
}
