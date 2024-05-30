package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/rosedblabs/rosedb/v2"
)

var db *rosedb.DB

var kvMapOldest = map[uint64]uint64{}
var kvMapLatest = map[uint64]uint64{}
var kvMapCount = 0

func OpenDb() {
	if db == nil {
		var err error
		// db, err = leveldb.OpenFile("./leveldb", nil)
		options := rosedb.DefaultOptions
		options.DirPath = "/tmp/rosedb_basic"
		db, err = rosedb.Open(options)
		if err != nil {
			panic(err)
		}
		start := time.Now()
		db.Merge(true)
		fmt.Println("Merge", time.Since(start))
	}
}

func KVWriteBehind() {

	OpenDb()
	batch := db.NewBatch(rosedb.DefaultBatchOptions)
	for key, value := range kvMapLatest {
		batch.Put(Uint64ToByteArray(key), Uint64ToByteArray(value))
	}
	_ = batch.Commit()

	kvMapLatest = map[uint64]uint64{}

}

var Uint64ToByteArray = func(i uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, i)
	return bs
}

func PutUint64(key uint64, value uint64) {
	kvMapLatest[key] = value
	kvMapCount++
	db.Put(Uint64ToByteArray(key), Uint64ToByteArray(value))
	// if kvMapCount >= 1_000_000 {
	// 	start := time.Now()
	// 	KVWriteBehind()
	// 	fmt.Println("KVWriteBehind", time.Since(start))
	// 	kvMapCount = 0
	// }
}

func GetUint64(key uint64) uint64 {
	if val, ok := kvMapLatest[key]; ok {
		return val
	}
	OpenDb()
	val, err := db.Get(Uint64ToByteArray(uint64(key)))
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(val)
}
