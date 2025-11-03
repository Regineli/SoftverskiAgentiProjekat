package model

import (
    "fmt"
)

// MusicModel represents a simple model with genre preference weights.
type MusicModel struct {
	Weights []float64 // e.g. [Rock, Pop, Jazz]
}

// NewModel initializes a model with random or zeroed weights.
func NewModel(numGenres int) *MusicModel {
	m := &MusicModel{Weights: make([]float64, numGenres)}
	for i := 0; i < numGenres; i++ {
		m.Weights[i] = 0.5 // neutral starting point
	}
	return m
}

// Print displays the model weights nicely.
func (m *MusicModel) Print() {
	fmt.Printf("🎧 Model Weights: ")
	for _, w := range m.Weights {
		fmt.Printf("%.2f ", w)
	}
	fmt.Println()
}
