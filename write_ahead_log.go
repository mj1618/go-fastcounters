package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/vmihailenco/msgpack/v5"
)

type WALEntry struct {
	Command        []byte
	CommandType    string
	SequenceNumber int64
	TraceId        int64
}

func GetCommand[K any](e WALEntry) K {
	var value K
	msgpack.Unmarshal(e.Command, &value)
	return value
}

var sequenceNumber atomic.Int64 = atomic.Int64{}
var traceId atomic.Int64 = atomic.Int64{}
var proposeWALEntryChannel chan WALEntry = make(chan WALEntry, 509600)
var readWALFile *os.File
var readWALDecoder *msgpack.Decoder
var writeWALFile *os.File
var responseWriterMap = sync.Map{}
var updateStateFn UpdateStateFunction = nil
var triggerReadWALFileChannel = make(chan int)

type UpdateStateFunction func(WALEntry)

func InitWriteAheadLog(fn UpdateStateFunction) {
	sequenceNumber.Store(0)
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
	go CommandLogReader()
	go CommandLogWriter()
}

func ProposeCommandToWAL(commandType string, command any) chan int {
	data, err := msgpack.Marshal(command)
	if err != nil {
		panic(err)
	}
	entry := WALEntry{Command: data, TraceId: traceId.Add(1), CommandType: commandType}
	proposeWALEntryChannel <- entry
	responseChannel := make(chan int)
	responseWriterMap.Store(entry.TraceId, responseChannel)
	return responseChannel
}

func ReadProposedWALEntries() []WALEntry {
	entries := make([]WALEntry, 0, 256)
	timer := time.NewTimer(time.Microsecond * 100)
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

func CommandLogReader() {
	for {
		n := <-triggerReadWALFileChannel
		entries := ReadCommandsFromWAL(n)
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
			triggerReadWALFileChannel <- len(entries)
		}
	}
}

func PreProcessCommands() {
	var i int64 = 0
	start := time.Now()
	for {
		commands := ReadCommandsFromWAL(5096)
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

	if i == 0 {
		fmt.Println("No commands found in WAL")
		fmt.Println()
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Read WAL in: \t%s\t seconds\n", humanize.Comma(int64(time.Since(start).Seconds())))
	fmt.Fprintf(w, "Read WAL at: \t%s\t msgs/sec\n", humanize.CommafWithDigits(float64(i)/time.Since(start).Seconds(), 0))
	fmt.Fprintf(w, "Commands read from file: \t%s\t\n", humanize.Comma(i))
	fmt.Fprintf(w, "Next sequence number: \t%s\t\n", humanize.Comma(sequenceNumber.Load()+1))
	w.Flush()
	fmt.Println()
	// fmt.Printf("Read WAL in %s\n", humanize.RelTime(start, time.Now(), "", ""))
	// fmt.Printf("Read WAL at %s msgs/sec\n", humanize.CommafWithDigits(float64(i)/time.Since(start).Seconds(), 0))
	// fmt.Printf("Total commands deserialised from file: %s\n", humanize.Comma(i))
	// fmt.Printf("Next sequence number: %s\n", humanize.Comma(sequenceNumber.Load()+1))
	// fmt.Println()
}

func ReadCommandsFromWAL(n int) (commands []WALEntry) {
	commands = make([]WALEntry, 0, n)
	for {
		command := WALEntry{}
		err := readWALDecoder.Decode(&command)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		commands = append(commands, command)
		if len(commands) == n {
			break
		}
	}
	return commands
}
