package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v3"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := badger.Open(badger.DefaultOptions(filepath.Join(home, "tasks.db")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 100; i++ {
		if err := db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(time.Now().String()), []byte("set at "+time.Now().String()))
		}); err != nil {
			log.Fatalln(err)
		}
	}

	if err := db.View(func(txn *badger.Txn) error {
		var it = txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			var item = it.Item()
			var v string
			if err := item.Value(func(val []byte) error {
				v = string(val)
				return nil
			}); err != nil {
				log.Fatalln(err)
			}
			log.Println(string(item.Key()), v)
		}
		return nil
	}); err != nil {
		log.Fatalln(err)
	}
}
