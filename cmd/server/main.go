package main

import (
	"log"

	"github.com/DedAzaMarks/ABS/internal/handler"
)

func main() {
	log.SetPrefix("server: ")
	log.Print("create")
	srv, err := handler.NewServer()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("run")
	log.Fatal(srv.Run())
}
