package model

import (
	"math"
)

// TrainOnUserData ažurira lokalni model korisnika.
// Ne koristi globalne mape jer svaki klijent ima svoj model.
func (m *MusicModel) TrainOnUserData(genreIndex int) {
	learningRate := 0.3

	for i := range m.Weights {
		if i == genreIndex {
			// povećaj preferencu za slušani žanr
			m.Weights[i] += learningRate * (1 - m.Weights[i])
		} else {
			// blago smanji ostale žanrove
			m.Weights[i] -= learningRate * m.Weights[i] * 0.2
		}

		// ograniči opseg između 0 i 1
		m.Weights[i] = math.Max(0, math.Min(1, m.Weights[i]))
	}
}
