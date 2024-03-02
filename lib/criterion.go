package lib

import (
	"fmt"
)

const (
	Benefit bool = true
	Cost    bool = false
)

type Criterion struct {
	weight         Evaluated
	typeOfCriteria bool
}

func (c *Criterion) set(value Evaluated, typeOF bool) {
	c.weight = value
	c.typeOfCriteria = typeOF
}

func (c *Criterion) String() string {
	s := "[" + fmt.Sprint(c.weight) + " "

	if c.typeOfCriteria == Benefit {
		s += "Benefit]"
	} else {
		s += "Cost]"
	}

	return s
}

func ChangeTypes(criteria []Criterion) []Criterion {
	newCriteria := make([]Criterion, len(criteria))
	for i := range newCriteria {
		newCriteria[i] = Criterion{criteria[i].weight, !criteria[i].typeOfCriteria}
	}
	return newCriteria
}
