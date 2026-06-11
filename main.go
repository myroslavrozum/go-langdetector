package main

import (
	"log"
	"os"
	"os/signal"

	db "go-langdetector/db"
)

func main() {
	c := make(chan os.Signal, 1)
	defer close(c)

	signal.Notify(c, os.Interrupt)

	database, err := db.InitDB("langdetector-badger-db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	go train(database)

	s := <-c
	log.Println("Got signal:", s)

}
