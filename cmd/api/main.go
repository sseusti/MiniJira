package main

import (
	"MiniJira/internal/httpapi"
	"MiniJira/internal/store/memory"
	"log"
	"net/http"
)

func main() {
	s := memory.NewStore()
	mux := httpapi.NewMux(s)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
