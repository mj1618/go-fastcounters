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
	if entry.CommandType == "MoveCommand" {
		countMoveCommands += 1
	} else if entry.CommandType == "IncrementCommand" {
		countIncrementCommands++
	} else {
		fmt.Println("Unknown command: ", entry)
	}
}

func GetCommandCounts() map[string]int {
	return map[string]int{
		"MoveCommands":      countMoveCommands,
		"IncrementCommands": countIncrementCommands,
	}
}
