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
// HandleUpload receives model updates from clients
func HandleUpload(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var update ModelUpdate
	_ = json.Unmarshal(body, &update)

	// proveri da li već postoji update od istog klijenta
	found := false
	for i, u := range allUpdates {
		if u.ClientID == update.ClientID {
			allUpdates[i] = update // zameni stari model novim
			found = true
			break
		}
	}

	if !found {
		allUpdates = append(allUpdates, update)
	}

	fmt.Printf("✅ Received update from Client %d: %v\n", update.ClientID, update.Weights)
	w.Write([]byte(fmt.Sprintf("Server stored model update for client %d\n", update.ClientID)))
}


// HandleAggregate performs FedAvg and returns global model
func HandleAggregate(w http.ResponseWriter, r *http.Request) {
    var models [][]float64

    // Sakupi sve modele od klijenata
    for _, u := range allUpdates {
        models = append(models, u.Weights)
    }

    if len(models) == 0 {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("⚠️ Nema dostupnih modela za agregaciju.\n"))
        return
    }

    // Izvrši federativnu agregaciju (FedAvg)
    global := FedAvg(models)

    fmt.Printf("🌍 Federated averaging završeno — globalni model: %v\n", global)

    // Očisti prethodne update-e nakon agregacije
    //allUpdates = nil

    // Vrati globalni model klijentu kao JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(global)
}



