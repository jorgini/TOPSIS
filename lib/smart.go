package lib

func (sm *SmartMatrix) CalcFinalScore() {
	for i := range sm.data {
		sm.finalScores[i] = sm.data[i].Sum()
	}
}
