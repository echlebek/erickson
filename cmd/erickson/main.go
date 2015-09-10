package main

import (
	"log"
	"net/http"
	"time"

	"github.com/echlebek/erickson/db"
	"github.com/echlebek/erickson/server"
)

func main() {
	db, err := db.NewBoltDB("my2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler := server.NewRootHandler(db, ".")

	s := http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
