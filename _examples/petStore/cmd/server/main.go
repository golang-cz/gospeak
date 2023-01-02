package main

import (
	"net/http"

	"github.com/golang-cz/gospeak/_examples/petStore/server"
)

func main() {
	api := &server.API{
		PetStore: map[int64]*server.Pet{},
	}

	handler := server.NewPetStoreServer(api)
	http.ListenAndServe(":8080", handler)
}
