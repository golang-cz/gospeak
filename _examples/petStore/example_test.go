package petStore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-cz/gospeak/_examples/petStore/client"
	"github.com/golang-cz/gospeak/_examples/petStore/server"
	"github.com/stretchr/testify/assert"
)

func TestInteroperability(t *testing.T) {
	srv := httptest.NewServer(server.NewPetStoreServer(&server.TestServer{}))
	defer srv.Close()

	api := client.NewPetStoreClient(srv.URL, &http.Client{})

	pets, err := api.ListPets(context.TODO())
	assert.NoError(t, err)
	assert.NotEmpty(t, pets)
}
