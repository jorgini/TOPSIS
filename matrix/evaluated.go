package matrix

type Evaluated interface {
	ConvertToNumbers() Number
	ConvertToInterval() Interval
	ConvertToT1FS(f Variants) *T1FS
	Weighted(weight Evaluated) Evaluated
	String() string
	DistanceNumber(other Evaluated, v Variants) (Number, error)
	DistanceInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error)
	Sum(other Evaluated) (Evaluated, error)
}
