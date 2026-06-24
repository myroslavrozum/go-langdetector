package main

import (
	"log"
	"os"
	"os/signal"

	"go-langdetector/db"
	"go-langdetector/trainer"
	"go-langdetector/webapp"
)

func main() {
	c := make(chan os.Signal, 1)
	logger := make(chan string)
	defer close(logger)
	defer close(c)

	signal.Notify(c, os.Interrupt)

	store, err := db.NewStore("data/langdetector-badger-db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	go trainer.Train(store, logger)
	go webapp.Run(store, logger)

	s := <-c
	log.Println("Got signal:", s)
}
