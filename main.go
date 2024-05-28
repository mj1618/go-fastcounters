package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/julienschmidt/httprouter"
	"github.com/vmihailenco/msgpack/v5"
)

var CommandTypes = map[string]int32{
	"MOVE":      1,
	"INCREMENT": 2,
	"DECREMENT": 3,
}

type Command struct {
	CommandType int32
	FromAddress int32
	ToAddress   int32
	Amount      int32
	RequestId   int64
}

var commandLogFile *os.File
var commandWriterChannel chan Command
var requestIdSeq atomic.Int64
var responseWriterMap = sync.Map{}

func main() {
	// count_file()
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

	// _, file = getDecoder("output.txt")
	defer file.Close()

	for {
		// err = decoder.Decode(&command)
		buf := make([]byte, 81)
		_, err := file.Read(buf)
		msgpack.Unmarshal(buf, &command)
		if err != nil {
			panic(err)
		}
		fmt.Println(command)
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
	_, file := getEncoder("output.txt")
	b := bufio.NewWriter(file)
	for {
		commands := ReadBatchOfCommands()
		// fmt.Println(len(commands))
		if len(commands) > 0 {
			for _, command := range commands {
				data, err := msgpack.Marshal(&command)

				if err != nil {
					panic(err)
				}
				b.Write(data)
				_ = data
			}

			// start := time.Now()
			// file.Sync()
			b.Flush()
			// fmt.Println("Time taken to write to file: ", time.Since(start))

			for _, command := range commands {
				responseChannel, ok := responseWriterMap.Load(command.RequestId)
				if ok {
					responseChannel.(chan int) <- 200
				}
			}
		}
	}
}

func ReadBatchOfCommands() []Command {
	commands := make([]Command, 0, 256)
	// timer := time.NewTimer(time.Microsecond * 1000)
	for {
		select {
		case command := <-commandWriterChannel:
			commands = append(commands, command)
			if len(commands) == 256 {
				return commands
			}
			// case <-timer.C:
			// 	return commands
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

func getEncoder(fileName string) (*gob.Encoder, *os.File) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	return gob.NewEncoder(file), file
}

func getDecoder(fileName string) (*gob.Decoder, *os.File) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	return gob.NewDecoder(file), file
}

// func writeFile(file *os.File, commands []Command) {

// 	for _, command := range commands {
// 		var data []byte
// 		var err error
// 		if data, err = proto.Marshal(&command); err != nil {
// 			log.Fatal(err)
// 		}
// 		_, err = file.Write(data)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	file.Sync()
// }
