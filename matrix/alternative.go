package matrix

import (
	"errors"
	"math"
)

type Alternative struct {
	grade           []Evaluated
	countOfCriteria int
}

func (a *Alternative) String() string {
	s := ""
	for i := 0; i < a.countOfCriteria; i++ {
		s += a.grade[i].String() + " "
	}
	return s
}

func (a *Alternative) FindDistance(to *Alternative) (Numbers, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return 0, errors.New("incomparable alternatives")
	}

	result := Numbers(0)

	for i := 0; i < a.countOfCriteria; i++ {
		if tmp, err := a.grade[i].DistanceTo(to.grade[i]); err != nil {
			return 0, errors.Join(err)
		} else {
			result += tmp
		}
	}

	result = Numbers(math.Sqrt(float64(result)))

	return result, nil
}
