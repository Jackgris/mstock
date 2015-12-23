package main

import (
	"log"

	"github.com/jackgris/mstock/server"
)

func main() {
	// only we create and launch the server
	log.Println("Create server")
	serve := server.NewServer()
	log.Println("Start server")
	serve.Run()
}
