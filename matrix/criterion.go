package matrix

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
