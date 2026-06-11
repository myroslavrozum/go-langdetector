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

	database, err := db.InitDB("data/langdetector-badger-db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	go train(database)
	go webapp.Run()

	s := <-c
	log.Println("Got signal:", s)

}
