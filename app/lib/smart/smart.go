package smart

func (sm *SmartMatrix) CalcFinalScore() {
	for i := range sm.Data {
		sm.FinalScores[i] = sm.Data[i].Sum()
	}
}
