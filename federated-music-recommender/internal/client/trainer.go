package client

import (
	"fmt"
	"federated-music-recommender/internal/model"
)

// EvaluateAndPrint poziva Evaluate() i ispisuje rezultate.
func EvaluateAndPrint(m *model.MusicModel, context string) {
	acc, f1 := m.Evaluate()
	if context != "" {
		fmt.Printf("📊 %s — Accuracy: %.2f | F1: %.2f\n", context, acc, f1)
	} else {
		fmt.Printf("📊 Evaluacija — Accuracy: %.2f | F1: %.2f\n", acc, f1)
	}
}
