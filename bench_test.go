package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// func TestBenchmark(t *testing.T) {
// 	InitWriteAheadLog(UpdateState)
// 	var wg sync.WaitGroup
// 	for i := 0; i <= 8; i++ {
// 		wg.Add(1)
// 		go (func() {
// 			defer wg.Done()
// 			for i := 0; i < 10000; i++ {
// 				responseChannel := ProposeCommandToWAL("MoveCommand", MoveCommand{FromAddress: 1, ToAddress: 2, Amount: 10})
// 				<-responseChannel
// 			}
// 		})()
// 	}
// 	wg.Wait()

// }

func runCommands(n int) {
	var wg sync.WaitGroup
	wg.Add(n)
	// var x atomic.Int32 // x is a counter
	for i := 0; i < n; i++ {
		go (func() {
			defer wg.Done()
			responseChannel := ProposeCommandToWAL("MoveCommand", MoveCommand{FromAddress: 1, ToAddress: 2, Amount: 10})
			<-responseChannel

			// fmt.Println("Done", x.Add(1))
		})()
	}
	// fmt.Println("Waiting for ", x, " commands to be processed")
	wg.Wait()
	// fmt.Println("Processed ", n, " commands")
}

func TestBench(t *testing.T) {
	InitWriteAheadLog(UpdateState)

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

// func TestBench(t *testing.T) {
// 	InitWriteAheadLog(UpdateState)

// 	var start = time.Now()
// 	var wg sync.WaitGroup
// 	n := 1000000
// 	m := 10
// 	wg.Add(n)
// 	var chanList = make([]chan int, m)

// 	for i := 0; i < n/m; i++ {
// 		for j := 0; j < m; j++ {
// 			chanList[i] = ProposeCommandToWAL("MoveCommand", MoveCommand{FromAddress: 1, ToAddress: 2, Amount: 10})
// 		}

// 		for x := 0; x < m; x++ {
// 			<-chanList[x]
// 		}
// 		fmt.Println("Processed ", i, " commands")
// 	}

// 	// wg.Wait()
// 	fmt.Println("Elapsed time: ", time.Since(start))
// 	fmt.Println("Commands processed: ", n)
// 	fmt.Println("Commands per second: ", float64(n)/time.Since(start).Seconds())
// }
