package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

var kvMap = map[uint64]uint64{}
var kvMapCount = 0

func OpenDb() {
	if db == nil {
		var err error
		db, err = leveldb.OpenFile("./leveldb", nil)
		if err != nil {
			panic(err)
		}
	}
}

func KVWriteBehind() {

	OpenDb()
	for i := 0; i < 100_000; i++ {
		for key, value := range kvMap {
			batch := new(leveldb.Batch)
			batch.Put(Uint64ToByteArray(key), Uint64ToByteArray(value))
			db.Write(batch, nil)

		}
	}
	kvMap = map[uint64]uint64{}

}

var Uint64ToByteArray = func(i uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, i)
	return bs
}

func PutUint64(key uint64, value uint64) {
	kvMap[key] = value
	kvMapCount++
	if kvMapCount >= 100_000 {
		start := time.Now()
		KVWriteBehind()
		fmt.Println("KVWriteBehind", time.Since(start))
		kvMapCount = 0
	}
}

func GetUint64(key uint64) uint64 {
	if val, ok := kvMap[key]; ok {
		return val
	}
	OpenDb()
	val, err := db.Get(Uint64ToByteArray(uint64(key)), nil)
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(val)
}
