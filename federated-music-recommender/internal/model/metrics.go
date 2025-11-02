package model

import (
	"fmt"
	"math/rand"
)

// Evaluate simulates evaluating the model and returns random accuracy/F1.
func (m *MusicModel) Evaluate() (float64, float64) {
	accuracy := 0.8 + rand.Float64()*0.2 // 0.8–1.0 range
	f1 := 0.75 + rand.Float64()*0.25     // 0.75–1.0 range
	fmt.Printf("📈 Evaluation - Accuracy: %.2f | F1: %.2f\n", accuracy, f1)
	return accuracy, f1
}
