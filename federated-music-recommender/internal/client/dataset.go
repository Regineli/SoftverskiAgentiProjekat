package client

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sort"
)

type Song struct {
	Name       string
	Artist     string
	Genre      string
	PlayCount  int
}

// LoadDataset učitava pesme sa brojem slušanja
func LoadDataset(path string) ([]Song, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("⚠️ Greška pri otvaranju fajla: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("⚠️ Greška pri čitanju CSV-a: %v", err)
	}

	var songs []Song
	genreMap := make(map[string]bool)

	for i, row := range rows {
		if i == 0 {
			continue // preskoči header
		}
		if len(row) < 4 {
			continue
		}

		playCount, _ := strconv.Atoi(row[3])

		s := Song{
			Name:      row[0],
			Artist:    row[1],
			Genre:     row[2],
			PlayCount: playCount,
		}
		songs = append(songs, s)
		genreMap[s.Genre] = true
	}

	var genres []string
	for g := range genreMap {
		genres = append(genres, g)
	}

	sort.Strings(genres)

	fmt.Printf("✅ Učitano %d pesama, %d žanrova.\n", len(songs), len(genres))
	return songs, genres, nil
}
