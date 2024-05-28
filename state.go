package main

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
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
	fmt.Println("Updating state with entry: ", entry)
	if entry.CommandType == "MoveCommand" {
		countMoveCommands += 1
		fmt.Println("MoveCommand count", countMoveCommands)
		var result MoveCommand
		mapstructure.Decode(entry.Command, &result)
		// fmt.Println("MoveCommand: ", result)
	} else if entry.CommandType == "IncrementCommand" {
		countIncrementCommands++
		var result IncrementCommand
		mapstructure.Decode(entry.Command, &result)
		// fmt.Println("IncrementCommand: ", result)
	} else {
		fmt.Println("Unknown command type", entry)

	}
}

func GetCommandCounts() map[string]int {
	return map[string]int{
		"MoveCommands":      countMoveCommands,
		"IncrementCommands": countIncrementCommands,
	}
}
