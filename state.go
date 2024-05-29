package main

import (
	"fmt"
)

type MoveCommand struct {
	FromAddress int64
	ToAddress   int64
	Amount      int64
}

type MoveAllCommand struct {
	FromAddress int64
	ToAddress   int64
}

type IncrementCommand struct {
	Address int64
	Amount  int64
}

type DecrementCommand struct {
	Address int64
	Amount  int64
}

var counters = map[int64]int64{}

var countCommands = 0

func UpdateState(entry WALEntry) {
	switch entry.CommandType {
	case "MoveCommand":
		countCommands++
		cmd := GetCommand[MoveCommand](entry)
		amount := min(counters[cmd.FromAddress], cmd.Amount)
		counters[cmd.FromAddress] -= amount
		counters[cmd.ToAddress] += amount

	case "MoveAllCommand":
		countCommands++
		cmd := GetCommand[MoveAllCommand](entry)
		counters[cmd.ToAddress] += cmd.FromAddress
		counters[cmd.FromAddress] = 0

	case "IncrementCommand":
		countCommands++
		cmd := GetCommand[IncrementCommand](entry)
		counters[cmd.Address] += cmd.Amount

	case "DecrementCommand":
		countCommands++
		cmd := GetCommand[DecrementCommand](entry)
		if cmd.Amount <= counters[cmd.Address] {
			counters[cmd.Address] -= cmd.Amount
			return
		}

	default:
		fmt.Println("Unknown command: ", entry)
	}
}

func GetCommandCounts() map[string]int {
	return map[string]int{
		"Commands": countCommands,
	}
}

func GetCounterState() map[int64]int64 {
	return counters
}
