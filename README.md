# GoSpeak - Go `interface{}` as your API <!-- omit in toc -->

**NOTICE: Under development. Seeking early user feedback.**

GoSpeak is a simple RPC framework, a lightweight alternative to [gRPC](https://grpc.io/) and [Twirp](https://twitchtv.github.io/twirp/docs/intro.html), where Go code is your protobuf.

```go
package proto

//go:webrpc golang -server -pkg=server -out=./server/server.gen.go
//go:webrpc golang -client -pkg=client -out=./client/example.gen.go
type ExampleAPI interface {
	Ping(context.Context, *Ping) (*Pong, error)
}
```

GoSpeak generates client/server code via [webrpc-gen](https://github.com/webrpc/webrpc) tool, which renders code from Go templates. The REST API routes and JSON payload are defined per [webrpc](https://github.com/webrpc/webrpc) specs and can be exported as OpenAPI 3.x (Swagger) API documentation.

| Server | | Client  |
|---|---|---|
| Go | <=> | [Go](https://github.com/webrpc/gen-golang) |
| Go | <=> | [TypeScript](https://github.com/webrpc/gen-typescript) |
| Go | <=> | [JavaScript](https://github.com/webrpc/gen-javascript) |
| Go | <=> | [Swagger codegen generators](https://github.com/webrpc/gen-openapi#generate-clientdocs-via-openapi-generator)|


# Quick example <!-- omit in toc -->

- [1. Write a Go interface](#1-write-a-go-interface)
- [2. Add webrpc targets](#2-add-webrpc-targets)
- [3. Generate code](#3-generate-code)
- [4. Mount and serve the API server](#4-mount-and-serve-the-api-server)
- [5. Implement the server business logic](#5-implement-the-server-business-logic)
- [6. Use the generated client](#6-use-the-generated-client)


## 1. Write a Go interface

This is your service definition (think of it as of `protobuf` file).

```go
package proto

import "context"

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
```

## 2. Add webrpc targets

Generate Go server and Go client:

```diff
+//go:webrpc golang -server -pkg=server -out=./server/server.gen.go
+//go:webrpc golang -client -pkg=client -out=./client/example.gen.go
 type PetStore interface {
 	GetPet(ctx context.Context, ID int64) (pet *Pet, err error)
 	ListPets(ctx context.Context) (pets []*Pet, err error)
 	CreatePet(ctx context.Context, new *Pet) (pet *Pet, err error)
 	UpdatePet(ctx context.Context, ID int64, update *Pet) (pet *Pet, err error)
 	DeletePet(ctx context.Context, ID int64) error
 }
```

Generate TypeScript client and OpenAPI 3.x (Swagger) docs too:

```diff
 //go:webrpc golang -server -pkg=server -out=./server/server.gen.go
 //go:webrpc golang -client -pkg=client -out=./client/example.gen.go
+//go:webrpc typescript -client -out=./client/exampleClient.gen.ts
+//go:webrpc openapi -out=./docs/exampleApi.gen.yaml -title=PetStoreAPI
 type PetStore interface {
 	GetPet(ctx context.Context, ID int64) (pet *Pet, err error)
 	ListPets(ctx context.Context) (pets []*Pet, err error)
 	CreatePet(ctx context.Context, new *Pet) (pet *Pet, err error)
 	UpdatePet(ctx context.Context, ID int64, update *Pet) (pet *Pet, err error)
 	DeletePet(ctx context.Context, ID int64) error
 }
```

## 3. Generate code

Install [gospeak](https://github.com/golang-cz/gospeak/releases) and generate the webrpc code.

```bash
$ gospeak ./proto/api.go
            PetStore => ./server/server.gen.go ✓
            PetStore => ./client/client.gen.go ✓
            PetStore => ./docs/videoApi.gen.yaml ✓
            PetStore => ./client/videoDashboardClient.gen.ts ✓
```

## 4. Mount and serve the API server

```go
// cmd/petstore/main.go
package main

import "./server"

func main() {
	api := &server.Server{} // implements PetStore interface{}

	handler := server.NewPetStoreServer(api)
	http.ListenAndServe(":8080", handler)
}
```

## 5. Implement the server business logic

The generated server code
- Handles incoming REST API requests
- Unmarshals JSON request into method argument(s)
- Calls your RPC method, ie. server.GetPet(ctx, petID) (*Pet, error)
- Marshals return argument(s) into JSON response

What's left is the business logic. Implement the interface methods:

```go
// server/server.go
package server

// Implements PetStore interface{}.
type Server struct {
	/* DB connection, config etc. */
}

```

```go
// server/user.go
package server

func (s *Server) GetUser(ctx context.Context, uid string) (user *User, err error) {
	user, err := s.DB.GetUser(ctx, uid)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, Errorf(ErrNotFound, "failed to find user(%v)", uid)
		}
		return nil, WrapError(ErrInternal, err, "failed to fetch user(%v)", uid)
	}

	return user, nil
}
```

See [source code](./_examples/petStore/server/pets.go)

## 6. Use the generated client

```go
package main

import "./client"

func main() {
	api := client.NewPetStoreClient(*serverUrl, &http.Client{})

	pets, err := api.ListPets(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pets)
}
```


## Enjoy! <!-- omit in toc -->

..and let us know what you think in [discussions](https://github.com/golang-cz/gospeak/discussions).

# Authors <!-- omit in toc -->
- [golang.cz](https://golang.cz)
- [VojtechVitek](https://github.com/VojtechVitek)

# License <!-- omit in toc -->

[MIT license](./LICENSE)
