package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"federated-music-recommender/internal/client"
	"federated-music-recommender/internal/model"

	"github.com/asynkron/protoactor-go/actor"
)

type UserActor struct {
    behavior actor.Behavior
    model    *model.MusicModel
    reader   *bufio.Scanner
    genres   []string // <── dodaj ovo
    userID   int
}

func NewUserActor(userID int) actor.Actor {
	songs, genres, err := client.LoadDataset("data/songs.csv")
	if err != nil {
		fmt.Println("⚠️  Greška pri učitavanju dataset-a:", err)
	}

	act := &UserActor{
		model:    model.NewModel(len(genres)),
		reader:   bufio.NewScanner(os.Stdin),
		behavior: actor.NewBehavior(),
		userID:   userID,
		genres:   genres, // <── čuvamo redosled
	}

	fmt.Printf("📊 Učitano %d pesama i %d žanrova.\n", len(songs), len(genres))
	act.behavior.Become(act.Menu)
	return act
}


func (a *UserActor) Receive(ctx actor.Context) {
	a.behavior.Receive(ctx)
}

// === STATE 1: MENU ===
func (a *UserActor) Menu(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *actor.Started:
		fmt.Printf("🚀 Klijent #%d pokrenut.\n", a.userID)
		a.showMenu()
		a.readInput(ctx)

	case string:
		switch msg {
		case "1":
			a.behavior.Become(a.ChoosingSong)
			ctx.Send(ctx.Self(), "choose")
		case "2":
			a.behavior.Become(a.SelectingGenre)
			ctx.Send(ctx.Self(), "genre")
		case "3":
			fmt.Println("📤 Slanje lokalnog modela serveru...")
			client.EvaluateAndPrint(a.model, fmt.Sprintf("Evaluacija pre slanja (Korisnik #%d)", a.userID))
			client.SendModelToServer(a.userID, a.model)
			a.backToMenu(ctx)
		case "4":
			fmt.Println("⬇️ Preuzimanje globalnog modela sa servera...")
			client.GetGlobalModel(a.model)
			client.EvaluateAndPrint(a.model, "Globalni model")
			a.backToMenu(ctx)
		case "exit":
			fmt.Println("👋 Izlaz iz aplikacije...")
			os.Exit(0)
		default:
			fmt.Println("⚠️ Nepoznata komanda.")
			a.showMenu()
			a.readInput(ctx)
		}
	}
}

// === STATE 2: CHOOSING SONG ===
func (a *UserActor) ChoosingSong(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case string:
		if msg == "choose" {
			client.RecommendSongs(a.model, a.genres)
			a.backToMenu(ctx)
		}
	}
}

// === STATE 3: SELECTING GENRE ===
func (a *UserActor) SelectingGenre(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case string:
		if msg == "genre" {
			client.SelectGenre(a.model, a.genres)
			a.backToMenu(ctx)
		}
	}
}

// === HELPER FUNKCIJE ===
func (a *UserActor) showMenu() {
	fmt.Println("\n==============================")
	fmt.Printf("👤 Korisnik #%d | 🎧 Model Weights po žanrovima:\n", a.userID)

	// Koristi žanrove koji su učitani u strukturi UserActor
	if len(a.genres) == 0 {
		fmt.Println("⚠️ Nema učitanih žanrova za korisnika.")
	} else {
		// Ispis imena žanrova u jednom redu
		for _, g := range a.genres {
			fmt.Printf("%-12s ", g)
		}
		fmt.Println()

		// Ispis težina u istom redosledu
		for _, w := range a.model.Weights {
			fmt.Printf("%-12.2f ", w)
		}
		fmt.Println()
	}

	fmt.Println("==============================")

	fmt.Println(`
=== 🎵 Meni ===
1 -> Pogledaj preporučene pesme
2 -> Izaberi žanr koji želiš da slušaš
3 -> Pošalji model serveru
4 -> Preuzmi globalni model
exit -> Zatvori klijenta
`)
}



func (a *UserActor) readInput(ctx actor.Context) {
	fmt.Print("👉 Izbor: ")
	if a.reader.Scan() {
		cmd := strings.TrimSpace(strings.ToLower(a.reader.Text()))
		ctx.Send(ctx.Self(), cmd)
	}
}

func (a *UserActor) backToMenu(ctx actor.Context) {
	a.behavior.Become(a.Menu)
	a.showMenu()
	a.readInput(ctx)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("⚠️  Koristi: go run ./cmd/client <userID>")
		return
	}

	userID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("⚠️  Nevažeći userID:", os.Args[1])
		return
	}

	system := actor.NewActorSystem()
	root := system.Root

	props := actor.PropsFromProducer(func() actor.Actor {
		return NewUserActor(userID)
	})

	root.Spawn(props)
	select {}
}
