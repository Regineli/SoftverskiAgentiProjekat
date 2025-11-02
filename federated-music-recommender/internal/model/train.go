package model

import (
	"math/rand"
	"time"
)

// TrainOnUserData simulates local training by slightly shifting weights
// toward the user's listening pattern.
func (m *MusicModel) TrainOnUserData(userID int) {
	rand.Seed(time.Now().UnixNano())

	// Example taste patterns per user
	userTastes := map[int][]float64{
		1: {1.0, 0.0, 0.0}, // User 1: Rock only
		2: {0.1, 1.0, 0.0}, // User 2: Pop
		3: {0.0, 0.9, 0.1}, // User 3: Pop + Jazz
	}

	taste := userTastes[userID]
	for i := range m.Weights {
		// simulate gradient descent towards user taste
		m.Weights[i] += 0.2 * (taste[i] - m.Weights[i])
		// add slight randomness
		m.Weights[i] += (rand.Float64() - 0.5) * 0.05
	}
}
