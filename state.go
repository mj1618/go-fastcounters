package main

import (
	"fmt"

	"github.com/mj1618/go-fastcounters/wal"
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

func UpdateState(entry wal.WALEntry, replaying bool) {

	switch entry.CommandType {
	case "MoveCommand":
		countCommands++
		cmd := wal.UnmarshalCommand[MoveCommand](entry)
		cmd.FromAddress = uint64(countCommands)
		cmd.ToAddress = uint64(countCommands) << 32
		if counters[cmd.FromAddress] >= cmd.Amount {
			counters[cmd.FromAddress] -= cmd.Amount
			counters[cmd.ToAddress] += cmd.Amount
		}

	case "MoveAllCommand":
		countCommands++
		cmd := wal.UnmarshalCommand[MoveAllCommand](entry)
		counters[cmd.ToAddress] += cmd.FromAddress
		counters[cmd.FromAddress] = 0

	case "IncrementCommand":
		countCommands++
		cmd := wal.UnmarshalCommand[IncrementCommand](entry)
		counters[cmd.Address] += cmd.Amount

	case "DecrementCommand":
		countCommands++
		cmd := wal.UnmarshalCommand[DecrementCommand](entry)
		if cmd.Amount <= counters[cmd.Address] {
			counters[cmd.Address] -= cmd.Amount
		}

	default:
		fmt.Println("Unknown command: ", entry)
	}

}

func UnmarshalCommandCounts() map[string]int {
	return map[string]int{
		"Commands": countCommands,
	}
}

func GetCounterState() any {
	return counters[1234]
}
