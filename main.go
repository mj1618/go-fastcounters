package main

import "github.com/mj1618/go-fastcounters/wal"

func main() {
	wal.InitWAL("", UpdateState)
	StartHttp()
}
