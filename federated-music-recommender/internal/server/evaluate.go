package server

import (
	"fmt"
	"math"
	"federated-music-recommender/internal/model"
)

// EvaluateGlobalModel računa koliko globalni model odgovara korisničkim modelima.
// Vraća prosečan Accuracy i F1-score svih klijenata.
func EvaluateGlobalModel(global []float64, clientModels [][]float64) (float64, float64) {
	if len(clientModels) == 0 || len(global) == 0 {
		fmt.Println("⚠️ Nema dovoljno modela za evaluaciju.")
		return 0, 0
	}

	totalAcc := 0.0
	totalF1 := 0.0
	count := 0

	for _, client := range clientModels {
		if len(client) != len(global) {
			fmt.Println("⚠️ Preskačem klijenta sa različitim brojem žanrova.")
			continue
		}

		// --- MSE (Mean Squared Error) ---
		mse := 0.0
		for i := range global {
			diff := global[i] - client[i]
			mse += diff * diff
		}
		mse /= float64(len(global))
		accuracy := math.Max(0, 1-mse)

		// --- Cosine similarity kao F1-score ---
		dot, normA, normB := 0.0, 0.0, 0.0
		for i := range global {
			dot += global[i] * client[i]
			normA += global[i] * global[i]
			normB += client[i] * client[i]
		}

		if normA == 0 || normB == 0 {
			continue
		}

		f1 := dot / (math.Sqrt(normA) * math.Sqrt(normB))

		totalAcc += accuracy
		totalF1 += f1
		count++
	}

	if count == 0 {
		return 0, 0
	}

	avgAcc := totalAcc / float64(count)
	avgF1 := totalF1 / float64(count)

	fmt.Printf("🌐 Global Evaluation — Avg Accuracy: %.2f | Avg F1: %.2f (iz %d klijenata)\n",
		avgAcc, avgF1, count)

	return avgAcc, avgF1
}


func TrainCentralizedModel(global []float64) *model.MusicModel {
	m := &model.MusicModel{Weights: make([]float64, len(global))}
	copy(m.Weights, global)

	// simulacija centralizovanog učenja — "brže konvergira"
	for i := range m.Weights {
		m.Weights[i] = (m.Weights[i] + 1) / 2
	}

	acc, f1 := m.Evaluate()
	fmt.Printf("🏛️ Centralized Training — Accuracy: %.2f | F1: %.2f\n", acc, f1)
	return m
}
