package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mj1618/go-fastcounters/wal"
	"github.com/vmihailenco/msgpack/v5"
)

func TestGenerics(t *testing.T) {
	mv := MoveCommand{FromAddress: 1, ToAddress: 2, Amount: 10}
	data, _ := msgpack.Marshal(mv)
	mv2 := wal.UnmarshalCommand[MoveCommand](wal.WALEntry{CommandData: data, CommandType: "MoveCommand"})
	fmt.Println(mv2)
	fmt.Println(reflect.TypeOf(mv2))
}

func TestInterfaceUnmarshal(t *testing.T) {
	mv := MoveCommand{FromAddress: 1, ToAddress: 2, Amount: 10}
	data, _ := msgpack.Marshal(mv)
	var mv2 interface{}
	msgpack.Unmarshal(data, &mv2)

	fmt.Println(mv2)
}
