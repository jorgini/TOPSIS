package lib

import (
	"fmt"
	"math"
	"reflect"
)

type Evaluated interface {
	ConvertToNumber() Number
	ConvertToInterval() Interval
	ConvertToT1FS(f Variants) *T1FS
	ConvertToAIFS(f Variants) *AIFS
	ConvertToIT2FS(f Variants) *IT2FS
	GetForm() Variants
	Weighted(weight Evaluated) Evaluated
	String() string
	DiffNumber(other Evaluated, v Variants) (Number, error)
	DiffInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error)
	Sum(other Evaluated) Evaluated
}

func Max[T Evaluated](a, b T) Evaluated {
	if reflect.TypeOf(a) == reflect.TypeOf(NumbersMin) {
		return Number(math.Max(float64(a.ConvertToNumber()), float64(b.ConvertToNumber())))
	}

	if reflect.TypeOf(a) == reflect.TypeOf(Interval{}) {
		s := Number(math.Max(float64(a.ConvertToInterval().Start), float64(b.ConvertToInterval().Start)))
		f := Number(math.Max(float64(a.ConvertToInterval().End), float64(b.ConvertToInterval().End)))
		return Interval{s, f}
	}

	if reflect.TypeOf(a) == reflect.TypeOf(&T1FS{}) {
		maxVert := make([]Number, len(a.ConvertToT1FS(Default).vert))
		for i := range maxVert {
			maxVert[i] = Number(math.Max(float64(a.ConvertToT1FS(Default).vert[i]),
				float64(b.ConvertToT1FS(Default).vert[i])))
		}
		return NewT1FS(maxVert...)
	}

	if reflect.TypeOf(a) == reflect.TypeOf(&AIFS{}) {
		maxVert := make([]Number, len(a.ConvertToAIFS(Default).vert))
		for i := range maxVert {
			maxVert[i] = Number(math.Max(float64(a.ConvertToAIFS(Default).vert[i]),
				float64(b.ConvertToAIFS(Default).vert[i])))
		}
		minPi := Number(math.Min(float64(a.ConvertToAIFS(Default).pi), float64(b.ConvertToAIFS(Default).pi)))
		return NewAIFS(minPi, maxVert...)
	}

	fmt.Println("Call deprecated method max")
	return nil
}

func Min[T Evaluated](a, b T) Evaluated {
	if reflect.TypeOf(a) == reflect.TypeOf(NumbersMin) {
		return Number(math.Min(float64(a.ConvertToNumber()), float64(b.ConvertToNumber())))
	}

	if reflect.TypeOf(a) == reflect.TypeOf(Interval{}) {
		s := Number(math.Min(float64(a.ConvertToInterval().Start), float64(b.ConvertToInterval().Start)))
		f := Number(math.Min(float64(a.ConvertToInterval().End), float64(b.ConvertToInterval().End)))
		return Interval{s, f}
	}

	if reflect.TypeOf(a) == reflect.TypeOf(&T1FS{}) {
		minVert := make([]Number, len(a.ConvertToT1FS(Default).vert))
		for i := range minVert {
			minVert[i] = Number(math.Min(float64(a.ConvertToT1FS(Default).vert[i]),
				float64(b.ConvertToT1FS(Default).vert[i])))
		}
		return NewT1FS(minVert...)
	}

	if reflect.TypeOf(a) == reflect.TypeOf(&AIFS{}) {
		minVert := make([]Number, len(a.ConvertToAIFS(Default).vert))
		for i := range minVert {
			minVert[i] = Number(math.Min(float64(a.ConvertToAIFS(Default).vert[i]),
				float64(b.ConvertToAIFS(Default).vert[i])))
		}
		maxPi := Number(math.Max(float64(a.ConvertToAIFS(Default).pi), float64(b.ConvertToAIFS(Default).pi)))
		return NewAIFS(maxPi, minVert...)
	}

	fmt.Println("Call deprecated method min")
	return nil
}

func HighType(a, b reflect.Type) reflect.Type {
	hasInterval := false
	hasT1FS := false
	hasAIFS := false
	hasIT2FS := false

	types := []reflect.Type{a, b}

	for _, t := range types {
		if t == reflect.TypeOf(Interval{}) {
			hasInterval = true
		}
		if t == reflect.TypeOf(&T1FS{}) {
			hasT1FS = true
		}
		if t == reflect.TypeOf(&AIFS{}) {
			hasAIFS = true
		}
		if t == reflect.TypeOf(&IT2FS{}) {
			hasIT2FS = true
		}
	}

	ret := reflect.TypeOf(NumbersMin)
	if hasIT2FS {
		ret = reflect.TypeOf(&IT2FS{})
	} else if hasAIFS {
		ret = reflect.TypeOf(&AIFS{})
	} else if hasT1FS {
		ret = reflect.TypeOf(&T1FS{})
	} else if hasInterval {
		ret = reflect.TypeOf(Interval{})
	}
	return ret
}
