package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"federated-music-recommender/internal/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type InputActor struct {
	MenuPID *actor.PID
}

func (i *InputActor) Receive(ctx actor.Context) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("👉 Izbor: ")
		if !reader.Scan() {
			break
		}
		cmd := strings.TrimSpace(reader.Text())
		ctx.Send(i.MenuPID, &messages.MenuCommand{Cmd: cmd})
	}
}
