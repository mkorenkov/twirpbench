package main

import (
	"log"
	"net/http"

	"github.com/mkorenkov/twirpbench/internal/twirpoptimized"
	"github.com/mkorenkov/twirpbench/internal/twirpoptimized/rpc/bloat"
)

const listenAddr = ":8082"

func main() {
	twirpHandler := bloat.NewBloatServer(&twirpoptimized.Server{})

	log.Printf("Listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, twirpHandler); err != nil {
		log.Fatal(err)
	}
}
