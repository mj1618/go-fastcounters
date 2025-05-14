package wal

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
	CommandData    []byte
	CommandType    string
	SequenceNumber int64
	TraceId        int64
}

func UnmarshalCommand[K any](e WALEntry) K {
	var value K
	msgpack.Unmarshal(e.CommandData, &value)
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
var triggerReadWALFileChannel = make(chan int, 256)

type UpdateStateFunction func(entry WALEntry, replaying bool)

func InitWAL(name string, fn UpdateStateFunction) {
	sequenceNumber.Store(0)

	updateStateFn = fn
	var err error

	readWALFile, err = os.OpenFile("./wal_"+name+".bin", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	readWALDecoder = msgpack.NewDecoder(readWALFile)

	writeWALFile, err = os.OpenFile("./wal_"+name+".bin", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	ReplayCommands()
	// go PeriodicallyFsync()
	go CommandLogReader()
	go CommandLogWriter()
}

// func PeriodicallyFsync() {
// 	for {
// 		time.Sleep(1 * time.Second)
// 		writeWALFile.Sync()
// 	}
// }

func ProposeCommandToWAL(commandType string, command any) chan int {
	data, err := msgpack.Marshal(command)
	if err != nil {
		panic(err)
	}
	entry := WALEntry{CommandData: data, TraceId: traceId.Add(1), CommandType: commandType}
	responseChannel := make(chan int, 1)
	responseWriterMap.Store(entry.TraceId, responseChannel)
	proposeWALEntryChannel <- entry
	return responseChannel
}

func ReadProposedWALEntries() []WALEntry {
	var entries = make([]WALEntry, 0, 256)
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
		// fmt.Println("Read ", len(entries), " entries from WAL")
		for _, entry := range entries {
			updateStateFn(entry, false)
			responseChannel, ok := responseWriterMap.Load(entry.TraceId)
			if ok {
				// fmt.Println("Found response channel for ", entry.TraceId)
				responseChannel.(chan int) <- 200
				// fmt.Println("Responded ", entry.TraceId)
				responseWriterMap.Delete(entry.TraceId)
			} else {
				fmt.Println("No response channel found for ", entry)
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
			// fmt.Println("Wrote ", len(entries), " entries to WAL")
			triggerReadWALFileChannel <- len(entries)
		}
	}
}

func ReplayCommands() {
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
				updateStateFn(command, true)
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
	fmt.Fprintf(w, "Read WAL in: \t%s\t seconds\n", humanize.CommafWithDigits((time.Since(start).Seconds()), 2))
	fmt.Fprintf(w, "Read WAL at: \t%s\t msgs/sec\n", humanize.CommafWithDigits(float64(i)/time.Since(start).Seconds(), 0))
	fmt.Fprintf(w, "Commands read from file: \t%s\t\n", humanize.Comma(i))
	fmt.Fprintf(w, "Next sequence number: \t%s\t\n", humanize.Comma(sequenceNumber.Load()+1))
	w.Flush()
	fmt.Println()
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
