package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mj1618/go-fastcounters/wal"
)

func runCommands(n int) {
	var wg sync.WaitGroup
	wg.Add(n)
	// var x atomic.Int32 // x is a counter
	for i := 0; i < n; i++ {
		go (func() {
			defer wg.Done()
			responseChannel := wal.ProposeCommandToWAL("MoveCommand", MoveCommand{FromAddress: 1, ToAddress: 2, Amount: 10})
			<-responseChannel

			// fmt.Println("Done", x.Add(1))
		})()
	}
	// fmt.Println("Waiting for ", x, " commands to be processed")
	wg.Wait()
	// fmt.Println("Processed ", n, " commands")
}

func TestBench(t *testing.T) {
	wal.InitWAL("testwal", UpdateState)

	var start = time.Now()
	n := 100_000_000
	batchSize := 100_000
	for i := 0; i < n; i += batchSize {
		runCommands(batchSize)
		// fmt.Println("Processed ", i+batchSize, " commands")
	}
	fmt.Println("Elapsed time: ", time.Since(start))
	fmt.Println("Commands processed: ", n)
	fmt.Println("Commands per second: ", float64(n)/time.Since(start).Seconds())
}
