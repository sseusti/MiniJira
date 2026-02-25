package main

import (
	"log"
	"net/http"
)

func main() {
	s := NewStore()
	mux := NewMux(s)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
