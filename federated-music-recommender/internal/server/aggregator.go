package server

import "fmt"

// FedAvg averages all clients' weight vectors to produce a new global model.
func FedAvg(models [][]float64) []float64 {
	if len(models) == 0 {
		return nil
	}

	numClients := len(models)
	numWeights := len(models[0])
	avg := make([]float64, numWeights)

	for _, weights := range models {
		for i := range weights {
			avg[i] += weights[i]
		}
	}

	for i := range avg {
		avg[i] /= float64(numClients)
	}

	fmt.Printf("🌍 FedAvg aggregated global model: %v\n", avg)
	return avg
}
