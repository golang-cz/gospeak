# GoSpeak - Generate REST APIs from Go code <!-- omit in toc -->

***NOTICE:** Under development. Seeking early user feedback.*

GoSpeak is a simple RPC framework, a lightweight alternative to [gRPC](https://grpc.io/) and [Twirp](https://twitchtv.github.io/twirp/docs/intro.html), where Go code is your protobuf.

```go
//go:generate github.com/golang-cz/gospeak/cmd/gospeak ./
package proto

//go:webrpc golang -server -pkg=server -out=./server/server.gen.go
//go:webrpc golang -client -pkg=client -out=./client/example.gen.go
type ExampleAPI interface {
	Ping(context.Context, *Ping) (*Pong, error)
}
```

Usage:

```
$ go get github.com/golang-cz/gospeak/cmd/gospeak@latest
$ go generate
            ExampleAPI => ./server/server.gen.go ✓
            ExampleAPI => ./client/example.gen.go ✓
```

## Language support <!-- omit in toc -->

GoSpeak uses [webrpc-gen](https://github.com/webrpc/webrpc) tool to generate REST API client & server code using Go templates. The API routes and JSON payload are defined per [webrpc](https://github.com/webrpc/webrpc) specs and can be exported to OpenAPI 3.x (Swagger) documentation.

| Server   | | Client                                                                                                               |
|----------|---|----------------------------------------------------------------------------------------------------------------------|
| Go 1.22+ | <=> | [Go 1.17+](https://github.com/webrpc/gen-golang)                                                                     |
| Go 1.22+ | <=> | [TypeScript](https://github.com/webrpc/gen-typescript)                                                               |
| Go 1.22+ | <=> | [JavaScript (ES6)](https://github.com/webrpc/gen-javascript)                                                         |
| Go 1.22+ | <=> | [OpenAPI 3+](https://github.com/webrpc/gen-openapi) (Swagger documentation)                                     |
| Go 1.22+ | <=> | Any OpenAPI client [code generator](https://github.com/webrpc/gen-openapi#generate-clientdocs-via-openapi-generator) |

# Quick example <!-- omit in toc -->

- [1. Define service API](#1-define-service-api)
- [2. Add target language directives](#2-add-target-language-directives)
- [3. Generate code](#3-generate-code)
- [4. Mount the API server](#4-mount-the-api-server)
- [5. Implement the server business logic](#5-implement-the-server-business-logic)
- [6. Use the generated client](#6-use-the-generated-client)
- [7. Test your API](#7-test-your-api)


## 1. Define service API

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

## 2. Add target language directives

Generate Go server and Go client code with `go:webrpc` directives:

```diff
+//go:webrpc golang -server -pkg=server -out=./server/server.gen.go
+//go:webrpc golang -client -pkg=client -out=./client/example.gen.go
 type PetStore interface {
```

Generate TypeScript client and OpenAPI 3.x (Swagger) documentation:

```diff
 //go:webrpc golang -server -pkg=server -out=./server/server.gen.go
 //go:webrpc golang -client -pkg=client -out=./client/example.gen.go
+//go:webrpc typescript -client -out=./client/exampleClient.gen.ts
+//go:webrpc openapi -out=./docs/exampleApi.gen.yaml -title=PetStoreAPI
 type PetStore interface {
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

*NOTE: Alternatively, you can `go get github.com/golang-cz/gospeak` as your dependency and run `go generate` against `//go:generate github.com/golang-cz/gospeak/cmd/gospeak .` directive.*

## 4. Mount the API server

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

The generated server code already
- handles incoming REST API requests
- unmarshals JSON request into method argument(s)
- calls your RPC method implementation, ie. `server.GetPet(ctx, 1)``
- marshals return argument(s) into a JSON response

What's left is the business logic. Implement the interface methods:

```go
// rpc/server.go
package rpc

type Server struct {
	/* DB connection, config etc. */
}
```

```go
// rpc/user.go
package rpc

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
	api := client.NewPetStoreClient("http://localhost:8080", http.DefaultClient)

	pets, err := api.ListPets(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pets)
}
```

## 7. Test your API

```go
package test

import (
	"testing"
	"./client"
)

func TestAPI(t *testing.T){
	api := client.NewPetStoreClient("http://localhost:8080", http.DefaultClient)

	pets, err := api.ListPets(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pets)
}
```

## Enjoy! <!-- omit in toc -->

..and let us know what you think in [discussions](https://github.com/golang-cz/gospeak/discussions).

# Authors <!-- omit in toc -->
- [golang.cz](https://golang.cz)
- [VojtechVitek](https://github.com/VojtechVitek)

# License <!-- omit in toc -->

[MIT license](./LICENSE)
