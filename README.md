# GoSpeak - Go `interface{}` as your API <!-- omit in toc -->

GoSpeak is a simple RPC framework, a lightweight alternative to [gRPC](https://grpc.io/) and [Twirp](https://twitchtv.github.io/twirp/docs/intro.html), where Go code is your protobuf.

```go
package schema

type ServiceDefinition interface {
	Ping(context.Context, *Ping) (*Pong, error)
}
```

GoSpeak generates REST API clients in multiple languages, OpenAPI 3.x (Swagger) documentation and Go server handler code. It's built on top of [webrpc](https://github.com/webrpc/webrpc) JSON schema protocol & code-generation suite, which uses Go templates to generate code.

| Server | | Client  |
|---|---|---|
| Go | <=> | [Go](https://github.com/webrpc/gen-golang) |
| Go | <=> | [TypeScript client](https://github.com/webrpc/gen-typescript) |
| Go | <=> | [JavaScript client](https://github.com/webrpc/gen-javascript) |
| Go | <=> | [Swagger codegen client(s)](https://github.com/swagger-api/swagger-codegen#overview)|


**NOTICE: Under development. We're seeking user feedback.**

# Quick example <!-- omit in toc -->

- [1. Define service API with Go `interface{}`](#1-define-service-api-with-go-interface)
- [2. Generate code](#2-generate-code)
	- [Generated server code (HTTP handlers)](#generated-server-code-http-handlers)
	- [Generated Go client](#generated-go-client)
	- [Generated OpenAPI 3.x (Swagger) documentation](#generated-openapi-3x-swagger-documentation)
- [4. Implement the API `interface{}` (server business logic)](#4-implement-the-api-interface-server-business-logic)


## 1. Define service API with Go `interface{}`

```go
package schema

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

## 2. Generate code

Install [gospeak](./releases) and generate your server code (HTTP handlers), strongly typed clients (Go/TypeScript) and documentation in OpenAPI 3.x (Swagger) API.

Generate webrpc `.json` schema only. Now you can use [webrpc-gen](https://github.com/webrpc/webrpc#getting-started) cli to generate code.
```bash
gospeak ./schema.api.go json -out ./schema.json
```

Or.. you can generate multiple targets directly from gospeak target:
```bash
gospeak ./schema/api.go \
  json -out ./schema.json \
  golang -server -pkg server -out ./server/server.gen.go \
  golang -client -pkg client -out ./client/client.gen.go \
  typescript -client -out ../frontend/src/client.gen.ts \
  openapi -out ./openapi.yaml
```

### Generated server code (HTTP handlers)

```go
/* generated server code */
package server

// - Handles incoming REST API requests
// - Unmarshals JSON request into method argument(s)
// - Calls your RPC method, ie. server.GetPet(ctx, petID) (*Pet, error)
// - Marshals return argument(s) into JSON response
func NewPetStoreServer(server PetStore) http.Handler {}
```

### Generated Go client

```go
// cmd/listpets/main.go
package main

import "./client"

var serverUrl = flag.String("serverUrl", "", "server URL")

func main() {
	api := client.NewPetStoreClient(*serverUrl, &http.Client{}) // generated client

	pets, err := api.ListPets(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pets)
}
```

### Generated OpenAPI 3.x (Swagger) documentation

```
/* TODO *? 

## 3. Mount and serve the API

```go
// cmd/petstore/main.go
package main

import "./server"

func main() {
	api := &server.API{} // implements API interface{}

	handler := server.NewPetStoreServer(api)
	http.ListenAndServe(":8080", handler)
}
```

## 4. Implement the API `interface{}` (server business logic)

```go
// server/user.go
package server

func (s *API) GetUser(ctx context.Context, uid string) (user *User, err error) {
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

## Enjoy! <!-- omit in toc -->

..and let us know what you think in [discussions](https://github.com/golang-cz/gospeak/discussions).

# Authors <!-- omit in toc -->
- [golang.cz](https://golang.cz)
- [VojtechVitek](https://github.com/VojtechVitek)

# License <!-- omit in toc -->

[MIT license](./LICENSE)
