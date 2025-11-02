package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"federated-music-recommender/internal/model"
)

type UserServer struct {
	mu     sync.Mutex
	Model  *model.MusicModel
	ClientID int
}

// Pokreće lokalni server za komunikaciju sa centralnim
func (s *UserServer) Start() {
	http.HandleFunc("/model", s.handleModel)
	http.HandleFunc("/update", s.handleUpdate)

	fmt.Println("🌐 Lokalni user server pokrenut na http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}

// GET /model — vraća trenutni model
func (s *UserServer) handleModel(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.Model)
}

// POST /update — prima novi globalni model
func (s *UserServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var global []float64
	if err := json.NewDecoder(r.Body).Decode(&global); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.Model.Weights = global
	s.mu.Unlock()

	fmt.Println("⬇️ Primljen novi globalni model:", global)
	w.WriteHeader(http.StatusOK)
}
