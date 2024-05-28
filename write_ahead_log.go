package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

var sequenceNumber atomic.Int64 = atomic.Int64{}
var proposeCommandChannel chan Command = make(chan Command, 509600)
var readWALFile *os.File
var writeWALFile *os.File
var responseWriterMap = sync.Map{}
var traceIdSeq atomic.Int64 = atomic.Int64{}
var updateStateFn UpdateStateFunction = nil

type UpdateStateFunction func(Command)

func InitWriteAheadLog(fn UpdateStateFunction) {
	sequenceNumber.Store(-1)
	updateStateFn = fn
	var err error
	readWALFile, err = os.OpenFile("./wal.bin", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	writeWALFile, err = os.OpenFile("./wal.bin", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	PreProcessCommands()
	go CommandLogWriter()
}

func ProposeCommandToWAL(command Command) chan int {
	command.TraceId = traceIdSeq.Add(1)
	proposeCommandChannel <- command
	responseChannel := make(chan int)
	responseWriterMap.Store(command.TraceId, responseChannel)
	return responseChannel
}

func ReadProposedWALCommands() []Command {
	commands := make([]Command, 0, 256)
	timer := time.NewTimer(time.Microsecond * 1000)
	for {
		select {
		case command := <-proposeCommandChannel:
			commands = append(commands, command)
			if len(commands) == 256 {
				return commands
			}
		case <-timer.C:
			return commands
		}
	}
}

func CommandLogWriter() {
	b := bufio.NewWriter(writeWALFile)
	for {
		commands := ReadProposedWALCommands()
		if len(commands) > 0 {
			for _, command := range commands {
				command.SequenceNumber = sequenceNumber.Add(1)
				data, err := msgpack.Marshal(&command)

				if err != nil {
					panic(err)
				}
				b.Write(data)
				_ = data
			}

			b.Flush()

			for _, command := range commands {
				updateStateFn(command)
				responseChannel, ok := responseWriterMap.Load(command.TraceId)
				if ok {
					responseChannel.(chan int) <- 200
					responseWriterMap.Delete(command.TraceId)
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
			}
			i++
		}
	}
	s, _ := readWALFile.Stat()
	fmt.Println("File size (bytes) / Command struct (bytes): ", s.Size()/int64(commandByteSize))
	fmt.Println("Total commands deserialised from file: ", i)
	fmt.Println("Next sequence number: ", sequenceNumber.Load()+1)
}

func ReadCommandsFromWAL() (commands []Command) {
	commands = make([]Command, 0, 256)
	for {
		buf := make([]byte, commandByteSize)
		_, err := readWALFile.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		var command Command
		msgpack.Unmarshal(buf, &command)
		commands = append(commands, command)
		if len(commands) == 256 {
			break
		}
	}
	return commands
}
