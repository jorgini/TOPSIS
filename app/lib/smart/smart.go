package smart

import (
	"sort"
	"sync"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func (sm *SmartMatrix) CalcFinalScore() {
	var wg sync.WaitGroup
	wg.Add(len(sm.Data))
	for i := range sm.Data {
		go func(i int) {
			defer wg.Done()
			sm.FinalScores[i] = sm.Data[i].Sum()
		}(i)
	}
	wg.Wait()
}

func (sm *SmartMatrix) RankedList(ranking v.Variants) matrix.RankedList {
	set := make([]eval.Rating, len(sm.FinalScores))
	ind := make([]int, len(sm.FinalScores))
	for i := range set {
		ind[i] = i
		set[i] = sm.FinalScores[i].CopyEval()
	}

	if ranking == v.Sengupta {
		sort.Slice(ind, func(i, j int) bool {
			l := set[ind[i]].ConvertToInterval()
			r := set[ind[j]].ConvertToInterval()
			return l.SenguptaGeq(r)
		})

		sort.Slice(set, func(i, j int) bool {
			l := set[i].ConvertToInterval()
			r := set[j].ConvertToInterval()
			return l.SenguptaGeq(r)
		})
	} else {
		sort.Slice(ind, func(i, j int) bool {
			return set[ind[i]].ConvertToNumber() > set[ind[j]].ConvertToNumber()
		})
		sort.Slice(set, func(i, j int) bool {
			return set[i].ConvertToNumber() > set[j].ConvertToNumber()
		})
	}
	return matrix.RankedList{Coeffs: set, Order: ind}
}
