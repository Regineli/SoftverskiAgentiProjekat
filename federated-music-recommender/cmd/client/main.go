package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"federated-music-recommender/internal/client"
	"federated-music-recommender/internal/model"

	"github.com/asynkron/protoactor-go/actor"
)

type UserActor struct {
	behavior actor.Behavior
	model    *model.MusicModel
	reader   *bufio.Scanner
}

func NewUserActor() actor.Actor {
	act := &UserActor{
		model:    model.NewModel(3),
		reader:   bufio.NewScanner(os.Stdin),
		behavior: actor.NewBehavior(),
	}
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
		fmt.Println("🚀 Klijent pokrenut.")
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
			client.EvaluateAndPrint(a.model, "Evaluacija pre slanja")
			client.SendModelToServer(1, a.model) // pošalji model
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
			client.RecommendSongs(a.model)
			a.backToMenu(ctx)
		}
	}
}

// === STATE 3: SELECTING GENRE ===
func (a *UserActor) SelectingGenre(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case string:
		if msg == "genre" {
			client.SelectGenre(a.model)
			a.backToMenu(ctx)
		}
	}
}

// === HELPER FUNKCIJE ===
func (a *UserActor) showMenu() {
	fmt.Println(`
=== 🎵 Meni ===
1 -> Pogledaj preporučene pesme
2 -> Izaberi žanr koji želiš da slušaš
3 -> Pošalji model serveru (simulacija)
4 -> Preuzmi globalni model (simulacija)
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
	system := actor.NewActorSystem()
	root := system.Root
	props := actor.PropsFromProducer(NewUserActor)
	root.Spawn(props)
	select {}
}
