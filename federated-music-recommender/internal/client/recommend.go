package client

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"federated-music-recommender/internal/model"
)

func RecommendSongs(m *model.MusicModel) {
	songs := []struct {
		name   string
		vector []float64
	}{
		{"Thunderstruck", []float64{1, 0, 0}},
		{"Bohemian Rhapsody", []float64{0.9, 0.1, 0}},
		{"Shape of You", []float64{0.1, 1, 0}},
		{"Uptown Funk", []float64{0.2, 0.8, 0}},
		{"Take Five", []float64{0, 0, 1}},
		{"Fly Me to the Moon", []float64{0, 0.1, 0.9}},
	}

	type scoredSong struct {
		name   string
		vector []float64
		score  float64
	}

	results := []scoredSong{}
	for _, s := range songs {
		score := 0.0
		for j := range s.vector {
			score += s.vector[j] * m.Weights[j]
		}
		results = append(results, scoredSong{s.name, s.vector, score})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	fmt.Println("\n🎵 Preporučene pesme:")
	for i, s := range results {
		fmt.Printf("%d. %-20s (score: %.2f)\n", i+1, s.name, s.score)
	}

	fmt.Print("👉 Unesi broj pesme (1–6) ili Enter za povratak: ")
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	input := strings.TrimSpace(reader.Text())
	if input == "" {
		return
	}

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(results) {
		fmt.Println("⚠️ Nevažeći izbor.")
		return
	}

	selected := results[index-1]
	fmt.Printf("🎧 Slušaš pesmu: %s\n", selected.name)
	m.TrainOnUserData(detectGenre(selected.vector))
	fmt.Println("📈 Model ažuriran.")
}

func detectGenre(v []float64) int {
	maxI := 0
	maxV := v[0]
	for i, val := range v {
		if val > maxV {
			maxI = i
			maxV = val
		}
	}
	return maxI + 1
}

func SelectGenre(m *model.MusicModel) {
	fmt.Println(`
🎧 Izaberi žanr:
1 -> Rock
2 -> Pop
3 -> Jazz
`)
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("👉 Unesi broj: ")
	reader.Scan()
	input := strings.TrimSpace(reader.Text())

	switch input {
	case "1":
		m.TrainOnUserData(1)
		fmt.Println("🎸 Slušao si Rock — model ažuriran.")
	case "2":
		m.TrainOnUserData(2)
		fmt.Println("🎤 Slušao si Pop — model ažuriran.")
	case "3":
		m.TrainOnUserData(3)
		fmt.Println("🎷 Slušao si Jazz — model ažuriran.")
	default:
		fmt.Println("⚠️ Nepoznat izbor.")
	}
}
