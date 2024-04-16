package eval

import (
	"encoding/json"
	"fmt"
	"math"
	v "webApp/lib/variables"
)

type Evaluated interface {
	GetType() string
	ConvertToNumber() Number
	ConvertToInterval() Interval
	ConvertToT1FS(f v.Variants) *T1FS
	ConvertToAIFS(f v.Variants) *AIFS
	ConvertToIT2FS(f v.Variants) *IT2FS
	GetForm() v.Variants
	Weighted(Weight Evaluated) Rating
	String() string
	DiffNumber(other Evaluated, variants v.Variants) (Number, error)
	DiffInterval(other Interval, typeOfCriterion bool, variants v.Variants) (Interval, error)
	Sum(other Evaluated) Rating
	CopyEval() Rating
	Equals(other Evaluated) bool
}

type Rating struct {
	Evaluated `json:"eval"`
}

func (r *Rating) IsNil() bool {
	return r.Evaluated == nil
}

func (r *Rating) MarshalJSON() ([]byte, error) {
	switch r.Evaluated.(type) {
	case Number:
		return json.Marshal(r.Evaluated.ConvertToNumber())
	case Interval:
		return json.Marshal(r.Evaluated.ConvertToInterval())
	case *T1FS:
		return json.Marshal(r.Evaluated.ConvertToT1FS(v.Default))
	case *IT2FS:
		return json.Marshal(r.Evaluated.ConvertToIT2FS(v.Default))
	case *AIFS:
		return json.Marshal(r.Evaluated.ConvertToAIFS(v.Default))
	case nil:
		return json.Marshal(nil)
	default:
		return nil, v.IncompatibleTypes
	}
}

func (r *Rating) UnmarshalJSON(data []byte) error {
	var t2 = IT2FS{Bottom: make([]Interval, 0), Upward: make([]Number, 0)}
	err := json.Unmarshal(data, &t2)
	if err == nil && len(t2.Bottom) > 0 {
		r.Evaluated = &t2
		return nil
	}

	var a = AIFS{T1FS: &T1FS{Vert: make([]Number, 0)}}
	err = json.Unmarshal(data, &a)
	if err == nil && len(a.Vert) > 0 {
		r.Evaluated = &a
		return nil
	}

	var t1 = T1FS{Vert: make([]Number, 0)}
	err = json.Unmarshal(data, &t1)
	if err == nil && len(t1.Vert) > 0 {
		r.Evaluated = &t1
		return nil
	}

	var i Interval
	err = json.Unmarshal(data, &i)
	if err == nil {
		r.Evaluated = i
		return nil
	}

	var n float64
	err = json.Unmarshal(data, &n)
	if err == nil {
		r.Evaluated = Number(n)
		return nil
	}

	return v.IncompatibleTypes
}

func Max[T Evaluated](a, b T) Rating {
	if a.GetType() == NumbersMin.GetType() {
		return Rating{Number(math.Max(float64(a.ConvertToNumber()), float64(b.ConvertToNumber())))}
	}

	if a.GetType() == (Interval{}).GetType() {
		s := Number(math.Max(float64(a.ConvertToInterval().Start), float64(b.ConvertToInterval().Start)))
		f := Number(math.Max(float64(a.ConvertToInterval().End), float64(b.ConvertToInterval().End)))
		return Rating{Interval{s, f}}
	}

	if a.GetType() == (&T1FS{}).GetType() {
		maxVert := make([]Number, len(a.ConvertToT1FS(v.Default).Vert))
		for i := range maxVert {
			maxVert[i] = Number(math.Max(float64(a.ConvertToT1FS(v.Default).Vert[i]),
				float64(b.ConvertToT1FS(v.Default).Vert[i])))
		}
		return Rating{NewT1FS(maxVert...)}
	}

	if a.GetType() == (&AIFS{}).GetType() {
		maxVert := make([]Number, len(a.ConvertToAIFS(v.Default).Vert))
		for i := range maxVert {
			maxVert[i] = Number(math.Max(float64(a.ConvertToAIFS(v.Default).Vert[i]),
				float64(b.ConvertToAIFS(v.Default).Vert[i])))
		}
		minPi := Number(math.Min(float64(a.ConvertToAIFS(v.Default).Pi), float64(b.ConvertToAIFS(v.Default).Pi)))
		return Rating{NewAIFS(minPi, maxVert...)}
	}

	fmt.Println("Call deprecated method max")
	return Rating{nil}
}

func Min[T Evaluated](a, b T) Rating {
	if a.GetType() == NumbersMin.GetType() {
		return Rating{Number(math.Min(float64(a.ConvertToNumber()), float64(b.ConvertToNumber())))}
	}

	if a.GetType() == (Interval{}).GetType() {
		s := Number(math.Min(float64(a.ConvertToInterval().Start), float64(b.ConvertToInterval().Start)))
		f := Number(math.Min(float64(a.ConvertToInterval().End), float64(b.ConvertToInterval().End)))
		return Rating{Interval{s, f}}
	}

	if a.GetType() == (&T1FS{}).GetType() {
		minVert := make([]Number, len(a.ConvertToT1FS(v.Default).Vert))
		for i := range minVert {
			minVert[i] = Number(math.Min(float64(a.ConvertToT1FS(v.Default).Vert[i]),
				float64(b.ConvertToT1FS(v.Default).Vert[i])))
		}
		return Rating{NewT1FS(minVert...)}
	}

	if a.GetType() == (&AIFS{}).GetType() {
		minVert := make([]Number, len(a.ConvertToAIFS(v.Default).Vert))
		for i := range minVert {
			minVert[i] = Number(math.Min(float64(a.ConvertToAIFS(v.Default).Vert[i]),
				float64(b.ConvertToAIFS(v.Default).Vert[i])))
		}
		maxPi := Number(math.Max(float64(a.ConvertToAIFS(v.Default).Pi), float64(b.ConvertToAIFS(v.Default).Pi)))
		return Rating{NewAIFS(maxPi, minVert...)}
	}

	fmt.Println("Call deprecated method min")
	return Rating{nil}
}

func HighType(a, b string) string {
	hasInterval := false
	hasT1FS := false
	hasAIFS := false
	hasIT2FS := false

	types := []string{a, b}

	for _, t := range types {
		if t == (Interval{}).GetType() {
			hasInterval = true
		}
		if t == (&T1FS{}).GetType() {
			hasT1FS = true
		}
		if t == (&AIFS{}).GetType() {
			hasAIFS = true
		}
		if t == (&IT2FS{}).GetType() {
			hasIT2FS = true
		}
	}

	ret := NumbersMin.GetType()
	if hasIT2FS {
		ret = (&IT2FS{}).GetType()
	} else if hasAIFS {
		ret = (&AIFS{}).GetType()
	} else if hasT1FS {
		ret = (&T1FS{}).GetType()
	} else if hasInterval {
		ret = (Interval{}).GetType()
	}
	return ret
}
