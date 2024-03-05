package lib

var (
	Bad       = 0
	Normal    = 1
	Good      = 2
	Excellent = 3
)

type LinguisticScale[T Evaluated] struct {
	Ratings []T
}

var DefaultNumberScale *LinguisticScale[Number] = &LinguisticScale[Number]{
	Ratings: []Number{1, 4, 6, 9},
}

var DefaultIntervalScale *LinguisticScale[Interval] = &LinguisticScale[Interval]{
	Ratings: []Interval{{1, 3}, {3, 5}, {5, 7}, {7, 9}},
}

var DefaultT1FSScale *LinguisticScale[*T1FS] = &LinguisticScale[*T1FS]{
	Ratings: []*T1FS{NewT1FS(0, 2, 3), NewT1FS(3, 4, 5), NewT1FS(5, 6, 7), NewT1FS(8, 9, 10)},
}
