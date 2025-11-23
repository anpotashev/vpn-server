package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anpotashev/vpn-server/internal/examplehandler"
)

func main() {
	port := 9999
	handler, err := examplehandler.NewExampleHandler()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", handler.ExampleHandler)
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	log.Fatal(server.ListenAndServe())
}
