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

// RecommendSongsLocal bira pesme iz dataset-a koje odgovaraju lokalnom modelu
func RecommendSongsLocal(m *model.MusicModel, allSongs []Song, genres []string) []Song {
	type scoredSong struct {
		Song
		Score float64
	}

	var scored []scoredSong

	// Mapiranje žanra na indeks
	genreIndex := make(map[string]int)
	for i, g := range genres {
		genreIndex[g] = i
	}

	// Izračunavanje score-a za svaku pesmu
	for _, s := range allSongs {
		if idx, ok := genreIndex[s.Genre]; ok && idx < len(m.Weights) {
			score := m.Weights[idx] * float64(s.PlayCount)
			scored = append(scored, scoredSong{Song: s, Score: score})
		}
	}

	// Sortiraj po score-u
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// Vrati top 10
	top := []Song{}
	for i := 0; i < len(scored) && i < 10; i++ {
		top = append(top, scored[i].Song)
	}

	return top
}

// RecommendSongs prikazuje personalizovane preporuke na osnovu CSV dataset-a
func RecommendSongs(m *model.MusicModel, genres []string) {
	songs, _, err := LoadDataset("data/songs.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	recommended := RecommendSongsLocal(m, songs, genres)
	if len(recommended) == 0 {
		fmt.Println("⚠️ Nema dostupnih pesama za preporuku.")
		return
	}

	fmt.Println("\n🎵 Personalizovane preporuke:")
	for i, s := range recommended {
		fmt.Printf("%2d. %-25s | %-15s | %-10s | %d slušanja\n",
			i+1, s.Name, s.Artist, s.Genre, s.PlayCount)
	}

	fmt.Print("👉 Unesi broj pesme (1–10) ili Enter za povratak: ")
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	input := strings.TrimSpace(reader.Text())
	if input == "" {
		return
	}

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(recommended) {
		fmt.Println("⚠️ Nevažeći izbor.")
		return
	}

	selected := recommended[index-1]
	fmt.Printf("🎧 Slušaš pesmu: %s — %s (%s)\n", selected.Name, selected.Artist, selected.Genre)

	// pronađi indeks žanra i ažuriraj lokalni model
	for i, g := range genres {
		if g == selected.Genre {
			m.TrainOnUserData(i)
			fmt.Printf("📈 Model ažuriran za žanr: %s\n", g)
			break
		}
	}
}

// SelectGenre omogućava ručni odabir žanra od strane korisnika
func SelectGenre(m *model.MusicModel, genres []string) {
	fmt.Println("\n🎧 Izaberi žanr:")

	for i, g := range genres {
		fmt.Printf("%d -> %s\n", i+1, g)
	}

	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("👉 Unesi broj: ")
	reader.Scan()
	input := strings.TrimSpace(reader.Text())

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(genres) {
		fmt.Println("⚠️ Nevažeći izbor.")
		return
	}

	selectedGenre := genres[index-1]
	m.TrainOnUserData(index - 1)
	fmt.Printf("🎧 Slušao si %s — lokalni model ažuriran.\n", selectedGenre)
}
