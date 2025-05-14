package main

import (
	"fmt"
	"time"

	// "github.com/mj1618/go-fastcounters/wal"
	"math/rand"

	"github.com/mj1618/go-fastcounters/wal"
	"github.com/notnil/chess"
)

func main() {
	wal.InitWAL("chess", UpdateChessState)
	InitChessState()
	ChessStartHttp()

	// wal.InitWAL("counters", UpdateState)
	// StartHttp()
}

func PlayChessGame() {
	game := chess.NewGame()

	// generate moves until game is over
	for game.Outcome() == chess.NoOutcome {
		startTime := time.Now()
		// select a random move
		moves := game.ValidMoves()
		move := moves[rand.Intn(len(moves))]
		game.Move(move)
		fmt.Println("2 Time taken: ", time.Since(startTime))
	}

	fmt.Println(game.Position().Board().Draw())
	fmt.Printf("Game completed. %s by %s.\n", game.Outcome(), game.Method())
	fmt.Println(game.String())
}
