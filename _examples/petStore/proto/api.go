package proto

import (
	"context"
	"time"

	"github.com/golang-cz/gospeak"
	"github.com/google/uuid"
)

//go:webrpc json -out=./petstore.gen.json
//go:webrpc debug -out=./petstore.debug.gen.txt
//go:webrpc golang -server -pkg=server -json=jsoniter -importTypesFrom=github.com/golang-cz/gospeak/_examples/petStore/proto -legacyErrors=true -out=./server/server.gen.go
//go:webrpc golang -client -pkg=client -json=jsoniter -out=./client/petstore.gen.go
//go:webrpc typescript -client -out=./petstore.gen.ts
//go:webrpc json -out=./petstore.gen.json
//go:webrpc openapi -out=./petstore.gen.yaml
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
	Tag     Tag `whatever`
	TagPtr  *Tag
	TagsPtr []*Tag

	Status Status `json:"status"`
}

type Tag struct {
	ID   int64
	Name string
}

// 0 = approved
// 1 = pending
// 2 = closed
// 3 = new
type Status gospeak.Enum[int]
