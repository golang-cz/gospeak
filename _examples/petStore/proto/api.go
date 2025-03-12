//go:generate go run github.com/golang-cz/gospeak/cmd/gospeak .
package proto

import (
	"context"
	"time"

	"github.com/golang-cz/gospeak/enum"
	"github.com/google/uuid"
)

//go:webrpc golang -server -pkg=proto -json=stdlib -types=false -out=./server.gen.go
//go:webrpc golang -client -pkg=client -json=stdlib -out=./client/petstore.gen.go
//go:webrpc typescript -client -out=./client/petstore.gen.ts
//go:webrpc openapi -out=./petstore.gen.yaml
//go:disabled json -out=./petstore.gen.json
//go:disabled debug -out=./petstore.debug.gen.txt
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

// approved = 0
// pending  = 1
// closed   = 2
// new      = 3
type Status enum.Int
