package petStore

import (
	"context"
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-cz/gospeak/_examples/petStore/proto"
	"github.com/golang-cz/gospeak/_examples/petStore/proto/client"
	"github.com/golang-cz/gospeak/_examples/petStore/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var serverUrl = flag.String("serverUrl", "", "server URL")

func TestPetStore(t *testing.T) {
	if *serverUrl == "" {
		// Run server, if not provided.
		api := &server.API{
			PetStore: map[int64]*proto.Pet{},
		}

		srv := httptest.NewServer(proto.NewPetStoreServer(api))
		defer srv.Close()

		*serverUrl = srv.URL
	}

	api := client.NewPetStoreClient(*serverUrl, &http.Client{})

	pets, err := api.ListPets(context.TODO())
	assert.NoError(t, err)
	assert.Empty(t, pets)

	pet, err := api.CreatePet(context.TODO(), &client.Pet{Name: "Daisy"})
	assert.NoError(t, err)
	require.NotNil(t, pet)

	_, err = api.GetPet(context.TODO(), pet.ID)
	assert.NoError(t, err)
	require.NotNil(t, pet)

	pets, err = api.ListPets(context.TODO())
	assert.NoError(t, err)
	assert.NotEmpty(t, pets)

	err = api.DeletePet(context.TODO(), pet.ID)
	assert.NoError(t, err)

	_, err = api.GetPet(context.TODO(), pet.ID)
	assert.Error(t, err)

	pets, err = api.ListPets(context.TODO())
	assert.NoError(t, err)
	assert.Empty(t, pets)
}
