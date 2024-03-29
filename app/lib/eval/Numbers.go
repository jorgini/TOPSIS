package eval

import (
	"fmt"
	"math"
	"reflect"
	v "webApp/lib/variables"
)

type Number float64

const NumbersMax Number = 2147483647

const NumbersMin Number = -2147483647

func (n Number) CopyEval() Rating {
	return Rating{n}
}

func (n Number) String() string {
	return fmt.Sprintf("%.2f", float64(n))
}

func (n Number) GetType() string {
	return reflect.TypeOf(n).String()
}

func (n Number) ConvertToNumber() Number {
	return n
}

func (n Number) ConvertToInterval() Interval {
	return Interval{n, n}
}

func (n Number) ConvertToT1FS(f v.Variants) *T1FS {
	if f == v.Default || f == v.Triangle {
		return NewT1FS(n, n, n)
	} else {
		return NewT1FS(n, n, n, n)
	}
}

func (n Number) ConvertToAIFS(f v.Variants) *AIFS {
	if f == v.Default || f == v.Triangle {
		return NewAIFS(0.0, n, n, n)
	} else {
		return NewAIFS(0.0, n, n, n, n)
	}
}

func (n Number) ConvertToIT2FS(f v.Variants) *IT2FS {
	if f == v.Default || f == v.Triangle {
		return NewIT2FS([]Interval{{n, n}, {n, n}}, []Number{n})
	} else {
		return NewIT2FS([]Interval{{n, n}, {n, n}}, []Number{n, n})
	}
}

func (n Number) GetForm() v.Variants {
	return v.None
}

func (n Number) Weighted(Weight Evaluated) Rating {
	if Weight.GetType() == NumbersMin.GetType() || Weight.GetType() == (Interval{}).GetType() {
		return Rating{n * Weight.ConvertToNumber()}
	}
	return Rating{nil}
}

func (n Number) DiffNumber(other Evaluated, variants v.Variants) (Number, error) {
	if other.GetType() != n.GetType() {
		return 0, v.IncompatibleTypes
	}

	if variants == v.SqrtDistance {
		return (n - other.ConvertToNumber()) * (n - other.ConvertToNumber()), nil
	} else if variants == v.CbrtDistance {
		return Number(math.Abs(math.Pow(float64(n-other.ConvertToNumber()), 3))), nil
	} else {
		return 0, v.InvalidCaseOfOperation
	}
}

func (n Number) DiffInterval(_ Interval, _ bool, _ v.Variants) (Interval, error) {
	return Interval{}, v.NoUsageMethod
}

func (n Number) Sum(other Evaluated) Rating {
	return Rating{n + other.ConvertToNumber()}
}
