package model

import "math/rand"

type MusicModel struct {
	Weights map[string]float64
}

func NewModel(genres []string) *MusicModel {
	weights := make(map[string]float64)
	for _, g := range genres {
		weights[g] = rand.Float64() * 0.5 // inicijalno random vrednosti
	}
	return &MusicModel{Weights: weights}
}

func (m *MusicModel) TrainOnGenre(genre string) {
	if _, exists := m.Weights[genre]; exists {
		m.Weights[genre] += 0.1 // svaka nova pesma jača preferencu
	}
}

func (m *MusicModel) GetPreferredGenres() []string {
	type kv struct {
		Key   string
		Value float64
	}
	var sorted []kv
	for k, v := range m.Weights {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	var result []string
	for _, kv := range sorted {
		result = append(result, kv.Key)
	}
	return result
}
