package matrix

import (
	"fmt"
	"webApp/lib/eval"
	v "webApp/lib/variables"
)

type Criterion struct {
	Weight         eval.Rating `json:"weight"`
	TypeOfCriteria bool        `json:"type_of_crit"`
}

func NewCriteria(size int) []Criterion {
	criteria := make([]Criterion, size)

	for i := range criteria {
		criteria[i] = Criterion{Weight: eval.Rating{Evaluated: eval.Number(0)}, TypeOfCriteria: v.Benefit}
	}
	return criteria
}

func CopyCriterion(c Criterion) Criterion {
	return Criterion{Weight: c.Weight.CopyEval(), TypeOfCriteria: c.TypeOfCriteria}
}

func (c *Criterion) set(value eval.Evaluated, typeOF bool) {
	c.Weight.Evaluated = value.CopyEval()
	c.TypeOfCriteria = typeOF
}

func (c *Criterion) String() string {
	s := "[" + fmt.Sprint(c.Weight) + " "

	if c.TypeOfCriteria == v.Benefit {
		s += "Benefit]"
	} else {
		s += "Cost]"
	}

	return s
}

func ChangeTypes(Criteria []Criterion) []Criterion {
	newCriteria := make([]Criterion, len(Criteria))
	for i := range newCriteria {
		newCriteria[i] = Criterion{Weight: Criteria[i].Weight, TypeOfCriteria: !Criteria[i].TypeOfCriteria}
	}
	return newCriteria
}

func GetHighType(Criteria []Criterion) string {
	highT := eval.NumbersMin.GetType()

	for _, c := range Criteria {
		highT = eval.HighType(highT, c.Weight.GetType())
	}
	return highT
}
