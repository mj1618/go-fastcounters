package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	InitWriteAheadLog(UpdateState)

	r := httprouter.New()
	r.GET("/", RootHandler)
	r.GET("/state", StateHandler)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func RootHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	responseChannel := ProposeCommandToWAL(Command{CommandType: CommandTypes["MOVE"], FromAddress: 1, ToAddress: 2, Amount: 10})
	result := <-responseChannel
	fmt.Fprintln(rw, "Result: ", result)
}

var currentState = 0

func UpdateState(command Command) {
	currentState++
}

func StateHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "State: ", currentState)
}
