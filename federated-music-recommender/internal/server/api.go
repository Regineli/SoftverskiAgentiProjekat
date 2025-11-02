package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ModelUpdate struct {
	ClientID int       `json:"client_id"`
	Weights  []float64 `json:"weights"`
}

var allUpdates []ModelUpdate

// HandleUpload receives model updates from clients
func HandleUpload(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var update ModelUpdate
	json.Unmarshal(body, &update)

	fmt.Printf("✅ Received update from Client %d: %v\n", update.ClientID, update.Weights)
	allUpdates = append(allUpdates, update)
	w.Write([]byte("Server received model update\n"))
}

// HandleAggregate performs FedAvg and returns global model
func HandleAggregate(w http.ResponseWriter, r *http.Request) {
	var models [][]float64
	for _, u := range allUpdates {
		models = append(models, u.Weights)
	}

	global := FedAvg(models)
	allUpdates = nil // clear after aggregation

	json.NewEncoder(w).Encode(global)
}
