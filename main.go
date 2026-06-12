package main

import (
	"log"
	"os"
	"os/signal"

	"go-langdetector/db"
	"go-langdetector/webapp"
)

func main() {
	c := make(chan os.Signal, 1)
	defer close(c)

	signal.Notify(c, os.Interrupt)

	store, err := db.NewStore("data/langdetector-badger-db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	go train(store)
	go webapp.Run()

	s := <-c
	log.Println("Got signal:", s)

}
