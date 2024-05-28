package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func getCommand[K any](data []byte) K {
	var x K
	msgpack.Unmarshal(data, &x)
	return x
}

func TestGenerics(t *testing.T) {
	mv := MoveCommand{FromAddress: 1, ToAddress: 2, Amount: 10}
	data, _ := msgpack.Marshal(mv)
	mv2 := getCommand[MoveCommand](data)
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
