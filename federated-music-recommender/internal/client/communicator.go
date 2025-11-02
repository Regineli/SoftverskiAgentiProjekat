package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"federated-music-recommender/internal/model"
)

type ModelUpdate struct {
	ClientID int       `json:"client_id"`
	Weights  []float64 `json:"weights"`
}

// SendModelToServer sends the trained model's weights to the central server.
func SendModelToServer(clientID int, m *model.MusicModel) {
	update := ModelUpdate{ClientID: clientID, Weights: m.Weights}
	data, _ := json.Marshal(update)

	resp, err := http.Post("http://localhost:8080/upload", "application/json", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("📤 Client %d sent model update: %v\n", clientID, m.Weights)
}

// ⬇️ Preuzima globalni model (HTTP GET /aggregate)
func GetGlobalModel(m *model.MusicModel) {
	resp, err := http.Get("http://localhost:8080/aggregate")
	if err != nil {
		fmt.Println("❌ Greška prilikom preuzimanja modela:", err)
		return
	}
	defer resp.Body.Close()

	var global []float64
	if err := json.NewDecoder(resp.Body).Decode(&global); err != nil {
		fmt.Println("❌ Neuspešno dekodiranje globalnog modela:", err)
		return
	}

	if len(global) > 0 {
		m.Weights = global
		fmt.Println("🌍 Preuzet globalni model sa servera:", global)
	} else {
		fmt.Println("⚠️ Nema dostupnih podataka za agregaciju.")
	}
}
