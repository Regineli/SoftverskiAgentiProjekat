package model

import (
	"fmt"
)

// MusicModel predstavlja jednostavan model sa žanrovskim težinama
// i internom memorijom o tome koliko je korisnik slušao svaki žanr.
type MusicModel struct {
	Weights    []float64 // normalizovane težine (zbir = 1)
	RawWeights []float64 // kumulativni brojevi slušanja (memorija)
}

// NewModel inicijalizuje model sa neutralnim vrednostima.
func NewModel(numGenres int) *MusicModel {
	m := &MusicModel{
		Weights:    make([]float64, numGenres),
		RawWeights: make([]float64, numGenres),
	}
	for i := 0; i < numGenres; i++ {
		m.Weights[i] = 1.0 / float64(numGenres) // ravnomerna početna distribucija
		m.RawWeights[i] = 1.0                   // početna vrednost da izbegnemo nulu u normalizaciji
	}
	return m
}

// Print prikazuje trenutne težine modela
func (m *MusicModel) Print() {
	fmt.Printf("🎧 Model Weights: ")
	for _, w := range m.Weights {
		fmt.Printf("%.2f ", w)
	}
	fmt.Println()
}
