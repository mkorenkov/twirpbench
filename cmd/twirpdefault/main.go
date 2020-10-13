package main

import (
	"log"
	"net/http"

	"github.com/mkorenkov/twirpbench/internal/twirpdefault"
	"github.com/mkorenkov/twirpbench/internal/twirpdefault/rpc/bloat"
)

const listenAddr = ":8081"

func main() {
	twirpHandler := bloat.NewBloatServer(&twirpdefault.Server{})

	log.Printf("Listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, twirpHandler); err != nil {
		log.Fatal(err)
	}
}
