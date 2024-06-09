package main

import (
	"fmt"
	"math/rand"

	"github.com/mj1618/go-fastcounters/wal"
	"github.com/notnil/chess"
)

var ngames = 10_000

var games map[int]*chess.Game = make(map[int]*chess.Game, ngames)
var gameMoveCounter map[int]int = make(map[int]int, ngames)

var gameId = 0

var gameCounter uint64 = 0
var gameCompleteCounter uint64 = 0
var moveCounter uint64 = 0

type ChessMoveCommand struct {
	move chess.Move
}

var blankGame = chess.NewGame()

var sampleMoves = []*chess.Move{}

func InitChessState() {
	for i := 0; i < ngames; i++ {
		games[i] = blankGame.Clone()
		gameMoveCounter[i] = 0
		gameCounter += 1
	}

	game := chess.NewGame()
	r := rand.New(rand.NewSource(123))
	for game.Outcome() == chess.NoOutcome {
		validMoves := game.ValidMoves()
		move := validMoves[r.Intn(len(validMoves))]
		game.Move(move)
		sampleMoves = append(sampleMoves, move)
	}
	fmt.Println("Sample moves: ", len(sampleMoves))

	testGame := chess.NewGame()
	moveI := 0
	for testGame.Outcome() == chess.NoOutcome {
		testGame.Move(sampleMoves[moveI])
		moveI++
	}
	fmt.Println("Test game complete: ", moveI)
}

func UpdateChessState(entry wal.WALEntry, replaying bool) {
	switch entry.CommandType {
	case "ChessMoveCommand":
		cmd := wal.UnmarshalCommand[ChessMoveCommand](entry)
		_ = cmd.move

		game := games[gameId]
		game.Move(sampleMoves[gameMoveCounter[gameId]])
		moveCounter++
		gameMoveCounter[gameId] = (gameMoveCounter[gameId] + 1) % len(sampleMoves)
		if game.Outcome() != chess.NoOutcome {
			games[gameId] = blankGame.Clone()
			gameMoveCounter[gameId] = 0
			gameCompleteCounter++
			gameCounter++
		}
		gameId = (gameId + 1) % ngames

	default:
		fmt.Println("Unknown command: ", entry)
	}

}

func GetChessState() map[string]uint64 {
	return map[string]uint64{
		"Games":          gameCounter,
		"Games Complete": gameCompleteCounter,
		"Moves":          moveCounter,
		"asfd":           uint64(gameMoveCounter[0]),
	}
}
