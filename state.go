package main

import (
	"fmt"
)

type MoveCommand struct {
	FromAddress int32
	ToAddress   int32
	Amount      int32
}

type IncrementCommand struct {
	Address int32
	Amount  int32
}

var countMoveCommands = 0
var countIncrementCommands = 0

func UpdateState(entry WALEntry) {
	switch entry.CommandType {
	case "MoveCommand":
		countMoveCommands++
		// GetCommand[MoveCommand](entry)
	case "IncrementCommand":
		countIncrementCommands++
		// GetCommand[IncrementCommand](entry)
	default:
		fmt.Println("Unknown command: ", entry)
	}
}

func GetCommandCounts() map[string]int {
	return map[string]int{
		"MoveCommands":      countMoveCommands,
		"IncrementCommands": countIncrementCommands,
	}
}
