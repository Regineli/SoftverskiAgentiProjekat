package model

import "math"

// TrainOnUserData ažurira lokalni model korisnika bez gubitka istorije
func (m *MusicModel) TrainOnUserData(genreIndex int) {
	learningRate := 0.3

	// ako nije inicijalizovano
	if m.RawWeights == nil {
		m.RawWeights = make([]float64, len(m.Weights))
	}

	// zabeleži "slušanje"
	m.RawWeights[genreIndex]++

	for i := range m.Weights {
		if i == genreIndex {
			m.Weights[i] += learningRate * (1 - m.Weights[i])
		} else {
			m.Weights[i] -= learningRate * m.Weights[i] * 0.2
		}
		m.Weights[i] = math.Max(0, math.Min(1, m.Weights[i]))
	}
}
