package main

import (
	"log"
	"net/http"

	"github.com/golang-cz/gospeak/_examples/petStore/server"
)

func main() {
	api := &server.API{
		PetStore: map[int64]*server.Pet{},
	}

	handler := server.NewPetStoreServer(api)

	log.Println("Serving PetStore API at :8080")
	http.ListenAndServe(":8080", handler)
}
