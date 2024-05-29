package main

import (
	"fmt"
)

type MoveCommand struct {
	FromAddress uint64
	ToAddress   uint64
	Amount      uint64
}

type MoveAllCommand struct {
	FromAddress uint64
	ToAddress   uint64
}

type IncrementCommand struct {
	Address uint64
	Amount  uint64
}

type DecrementCommand struct {
	Address uint64
	Amount  uint64
}

var counters = map[uint64]uint64{}

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

func GetCounterState() map[uint64]uint64 {
	return counters
}
