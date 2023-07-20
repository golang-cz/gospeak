package proto

import (
	"context"
	"time"

	"github.com/google/uuid"
)

//go:webrpc json -out=./petstore.gen.json
//go:webrpc /Users/vojtechvitek/webrpc/gen-golang -server -importTypesFrom=github.com/golang-cz/gospeak/_examples/petStore/proto -legacyErrors=true -pkg=server -out=./server/server.gen.go
//go:xxx /Users/vojtechvitek/webrpc/gen-golang -server -legacyErrors=true -pkg=server -out=./server/server.gen.go
//go:webrpc /Users/vojtechvitek/webrpc/gen-golang -client -pkg=client -out=./client/petstore.gen.go
//go:xxx typescript -client -out=./petstore.gen.ts
//go:xxx json -out=./petstore.gen.json
//go:xxx openapi -out=./petstore.gen.yaml
type PetStore interface {
	GetPet(ctx context.Context, ID int64) (pet *Pet, err error)
	ListPets(ctx context.Context) (pets []*Pet, err error)
	CreatePet(ctx context.Context, new *Pet) (pet *Pet, err error)
	UpdatePet(ctx context.Context, ID int64, update *Pet) (pet *Pet, err error)
	DeletePet(ctx context.Context, ID int64) error
}

type Pet struct {
	ID        int64      `json:"id,string"`
	UUID      uuid.UUID  `json:"uuid,string"`
	Name      string     `json:"name"`
	Available bool       `json:"available"`
	PhotoURLs []string   `json:"photoUrls"`
	Tags      []Tag      `json:"tags"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`

	// Test
	Tag     Tag `tag`
	TagPtr  *Tag
	TagsPtr []*Tag
}

type Tag struct {
	ID   int64
	Name string
}
