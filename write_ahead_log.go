package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type WALEntry struct {
	Command        any
	CommandType    string
	SequenceNumber int64
	TraceId        int64
}

var sequenceNumber atomic.Int64 = atomic.Int64{}
var traceId atomic.Int64 = atomic.Int64{}
var proposeWALEntryChannel chan WALEntry = make(chan WALEntry, 509600)
var readWALFile *os.File
var readWALDecoder *msgpack.Decoder
var writeWALFile *os.File
var responseWriterMap = sync.Map{}
var updateStateFn UpdateStateFunction = nil

type UpdateStateFunction func(WALEntry)

func InitWriteAheadLog(fn UpdateStateFunction) {
	sequenceNumber.Store(-1)
	updateStateFn = fn
	var err error

	readWALFile, err = os.OpenFile("./wal.bin", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	readWALDecoder = msgpack.NewDecoder(readWALFile)

	writeWALFile, err = os.OpenFile("./wal.bin", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	PreProcessCommands()
	go CommandLogWriter()
}

func ProposeCommandToWAL(commandType string, command any) chan int {
	entry := WALEntry{Command: command, TraceId: traceId.Add(1), CommandType: commandType}
	proposeWALEntryChannel <- entry
	responseChannel := make(chan int)
	responseWriterMap.Store(entry.TraceId, responseChannel)
	return responseChannel
}

func ReadProposedWALEntries() []WALEntry {
	entries := make([]WALEntry, 0, 256)
	timer := time.NewTimer(time.Microsecond * 1000)
	for {
		select {
		case entry := <-proposeWALEntryChannel:
			entries = append(entries, entry)
			if len(entries) == 256 {
				return entries
			}
		case <-timer.C:
			return entries
		}
	}
}

func CommandLogWriter() {
	b := bufio.NewWriter(writeWALFile)
	for {
		entries := ReadProposedWALEntries()
		if len(entries) > 0 {
			for _, entry := range entries {
				entry.SequenceNumber = sequenceNumber.Add(1)
				data, err := msgpack.Marshal(&entry)

				if err != nil {
					panic(err)
				}
				b.Write(data)
				_ = data
			}

			b.Flush()

			for _, entry := range entries {
				updateStateFn(entry)
				responseChannel, ok := responseWriterMap.Load(entry.TraceId)
				if ok {
					responseChannel.(chan int) <- 200
					responseWriterMap.Delete(entry.TraceId)
				}
			}
		}
	}
}

func PreProcessCommands() {
	i := 0
	for {
		commands := ReadCommandsFromWAL()
		if len(commands) == 0 {
			break
		}
		for _, command := range commands {
			if command.SequenceNumber > sequenceNumber.Load() {
				sequenceNumber.Store(command.SequenceNumber)
				updateStateFn(command)
			} else {
				panic(fmt.Sprintf("Sequence number mismatch %v %d", command, sequenceNumber.Load()))
			}
			i++
		}
	}

	fmt.Println("Total commands deserialised from file: ", i)
	fmt.Println("Next sequence number: ", sequenceNumber.Load()+1)
}

func ReadCommandsFromWAL() (commands []WALEntry) {
	commands = make([]WALEntry, 0, 256)
	for {
		command := WALEntry{}
		err := readWALDecoder.Decode(&command)

		// buf := make([]byte, commandByteSize)
		// _, err := readWALFile.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		commands = append(commands, command)
		if len(commands) == 256 {
			break
		}
		// var command Command
		// msgpack.Unmarshal(buf, &command)

	}
	return commands
}

// func ReadCommandsFromWAL() (commands []Command) {
// 	commands = make([]Command, 0, 256)
// 	for {
// 		buf := make([]byte, commandByteSize)
// 		_, err := readWALFile.Read(buf)
// 		if err == io.EOF {
// 			break
// 		} else if err != nil {
// 			panic(err)
// 		}
// 		var command Command
// 		msgpack.Unmarshal(buf, &command)
// 		commands = append(commands, command)
// 		if len(commands) == 256 {
// 			break
// 		}
// 	}
// 	return commands
// }
