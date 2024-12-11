package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SpectralJager/kv"
)

func main() {
	// open local database
	db, err := kv.Open("example/example.kv")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// write some data
	kv.Post(db, "name", "Daniel")
	kv.Post(db, "nickname", "SpectralJager")
	kv.Post(db, "year", 2024)
	kv.Post(db, "time", time.Now())

	// get available keys
	keys := kv.Head(db)
	for _, key := range keys {
		fmt.Println(key)
	}

	// delete some data
	kv.Delete(db, "year")

	// update some data
	kv.Put(db, "time", time.Now())

	// get some data
	name, _ := kv.Get[string](db, "name")
	nickname, _ := kv.Get[string](db, "nickname")
	time, _ := kv.Get[time.Time](db, "time")
	fmt.Printf("[%s] %s -> @%s\n", time.String(), name, nickname)
}
