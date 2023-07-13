package proto

import "context"

//go:webrpc golang@v0.10.0 -server -pkg=server -out=./server/server.gen.go
//go:webrpc golang@v0.10.0 -client -pkg=client -out=./client/petstore.gen.go
//go:webrpc typescript@v0.10.0 -client -out=./petstore.gen.ts
//go:webrpc json -out=./petstore.gen.json
//go:webrpc openapi@v0.10.0 -out=./petstore.gen.yaml
type PetStore interface {
	GetPet(ctx context.Context, ID int64) (pet *Pet, err error)
	ListPets(ctx context.Context) (pets []*Pet, err error)
	CreatePet(ctx context.Context, new *Pet) (pet *Pet, err error)
	UpdatePet(ctx context.Context, ID int64, update *Pet) (pet *Pet, err error)
	DeletePet(ctx context.Context, ID int64) error
}

type Pet struct {
	ID        int64
	Name      string
	Available bool
	PhotoURLs []string
	Tags      []Tag
}

type Tag struct {
	ID   int64
	Name string
}
