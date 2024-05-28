package main

import "github.com/vmihailenco/msgpack/v5"

var CommandTypes = map[string]int32{
	"MOVE":      1,
	"INCREMENT": 2,
	"DECREMENT": 3,
}

type Command struct {
	CommandType    int32
	FromAddress    int32
	ToAddress      int32
	Amount         int32
	TraceId        int64
	SequenceNumber int64
}

var commandByteSize int = getCommandByteSize()

func getCommandByteSize() int {
	command := Command{}
	data, err := msgpack.Marshal(&command)
	if err != nil {
		panic(err)
	}
	return len(data)
}
