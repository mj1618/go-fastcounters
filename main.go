package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/protobuf/proto"
)

var CommandTypes = map[string]int32{
	"MOVE":      1,
	"INCREMENT": 2,
	"DECREMENT": 3,
}

var commandLogFile *os.File
var commandWriterChannel chan Command
var requestIdSeq atomic.Int64
var responseWriterMap = sync.Map{}

func main() {
	count_file()
	commandLogFile = openLogFile()
	defer closeLogFile(commandLogFile)
	commandWriterChannel = make(chan Command, 256)
	go CommandLogWriter()

	r := httprouter.New()
	r.GET("/", RootHandler)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func count_file() {
	// open output file
	file, err := os.OpenFile("output.txt", os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var command Command
	i := 0

	buf := make([]byte, command)

	for {
		_, err = file.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}
		err = proto.Unmarshal(buf, &command)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(command.RequestId)
		i++
	}
	fmt.Println("Total commands: ", i)
	panic("done")

}

func RootHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	requestId := requestIdSeq.Add(1)
	responseWriterMap.Store(requestId, make(chan int))
	commandWriterChannel <- Command{CommandType: CommandTypes["MOVE"], FromAddress: 1, ToAddress: 2, Amount: 10, RequestId: requestId}
	responseChannel, ok := responseWriterMap.Load(requestId)
	if !ok {
		fmt.Fprintln(rw, "Error: Request Id not found")
		return
	}
	result := <-responseChannel.(chan int)
	fmt.Fprintln(rw, "Result: ", result)
	responseWriterMap.Delete(requestId)
}

func CommandLogWriter() {
	for {
		commands := ReadBatchOfCommands()
		if len(commands) > 0 {
			// fmt.Println("batched commands: ", len(commands))
			// writeChunk(EncodeCommands(commands))
			writeFile(commandLogFile, commands)
			for _, command := range commands {
				responseChannel, ok := responseWriterMap.Load(command.RequestId)
				if ok {
					responseChannel.(chan int) <- 200
				}
			}
		}
	}
}

func EncodeCommands(commands []Command) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(commands); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func DecodeCommand(buf bytes.Buffer) *Command {
	var command Command
	enc := gob.NewDecoder(&buf)
	err := enc.Decode(&command)
	if err != nil {
		panic(err)
	}
	return &command
}

func ReadBatchOfCommands() []Command {
	commands := make([]Command, 0, 5096)
	timer := time.NewTimer(time.Microsecond * 200)
	for {
		select {
		case command := <-commandWriterChannel:
			commands = append(commands, command)
			if len(commands) == 5096 {
				return commands
			}
		case <-timer.C:
			return commands
		}
	}
}

func openLogFile() (fo *os.File) {
	// open output file
	fo, err := os.Create("./output.txt")
	if err != nil {
		panic(err)
	}
	return fo
}

func closeLogFile(fo *os.File) {
	// close fo on exit and check for its returned error
	if err := fo.Close(); err != nil {
		panic(err)
	}
}

func writeFile(file *os.File, commands []Command) {

	for _, command := range commands {
		var data []byte
		var err error
		if data, err = proto.Marshal(&command); err != nil {
			log.Fatal(err)
		}
		_, err = file.Write(data)
		if err != nil {
			log.Fatal(err)
		}
	}
	file.Sync()
}
