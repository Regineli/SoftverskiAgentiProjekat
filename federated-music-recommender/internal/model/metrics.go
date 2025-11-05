package model

import (
	"fmt"
	"math"
)

// Evaluate poredi trenutne težine modela sa istorijom slušanja (RawWeights)
// i vraća metrike tačnosti (accuracy) i F1-score.
func (m *MusicModel) Evaluate() (float64, float64) {
	fmt.Println(len(m.RawWeights))
	if len(m.RawWeights) == 0 || len(m.RawWeights) != len(m.Weights) {
		fmt.Println("⚠️ Model nema dovoljno podataka za evaluaciju.")
		return 0, 0
	}

	// Normalizuj RawWeights (da zbir bude 1)
	sum := 0.0
	for _, v := range m.RawWeights {
		sum += v
	}
	trueDist := make([]float64, len(m.RawWeights))
	for i, v := range m.RawWeights {
		trueDist[i] = v / sum
	}

	// Izračunaj MSE grešku između modela i stvarne distribucije
	mse := 0.0
	for i := range m.Weights {
		diff := m.Weights[i] - trueDist[i]
		mse += diff * diff
	}
	mse /= float64(len(m.Weights))

	// Tačnost = 1 - MSE (bliže 1 znači bolja predikcija)
	accuracy := math.Max(0, 1-mse)

	// Izračunaj F1-score po formuli baziranoj na cosine similarity
	dot, normA, normB := 0.0, 0.0, 0.0
	for i := range m.Weights {
		dot += m.Weights[i] * trueDist[i]
		normA += m.Weights[i] * m.Weights[i]
		normB += trueDist[i] * trueDist[i]
	}
	if normA == 0 || normB == 0 {
		return 0, 0
	}
	cosineSim := dot / (math.Sqrt(normA) * math.Sqrt(normB))

	// Pretvori cosineSim u F1-score (0–1)
	f1 := cosineSim

	fmt.Printf("📈 Evaluation - Accuracy: %.2f | F1: %.2f\n", accuracy, f1)
	return accuracy, f1
}
