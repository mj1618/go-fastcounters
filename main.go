package main

func main() {
	InitWriteAheadLog(UpdateState)
	StartHttp()
}
