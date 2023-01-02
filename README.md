# GoSpeak - Go `interface{}` as your API

**NOTICE: Not stable. GoSpeak is under active development.**

What if Go `interface{}` was your schema for service-to-service communication? What if you could generate REST API server code, documentation and strongly typed clients in Go/TypesScript/JavaScript in seconds? What if you could use Go channels over network easily?

Introducing **GoSpeak**, a lightweight JSON alternative to gRPC and Twirp, where Go `interface{}` is your protobuf schema. GoSpeak is built on top of [webrpc](https://github.com/webrpc/webrpc) JSON protocol & code-generation suite.

## Example

1. Define your API schema with Go `interface{}`
2. Generate code (API handlers, Go/TypeScript clients, API docs)
3. Mount and serve the API
4. Implement the `interface{}` (server business logic)

### 1. Define your API schema with Go `interface{}`

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

### 2. Generate code (server handlers, Go/TS clients, API docs)

Install [gospeak](./releases) and generate your server code (HTTP handlers), strongly typed clients (Go/TypeScript) and documentation in OpenAPI 3.x (Swagger) API.

```bash
#!/bin/bash

gospeak ./schema/api.go \
  golang -server -pkg server -out ./server/server.gen.go \
  golang -client -pkg client -out ./client/client.gen.go \
  typescript -client -out ../frontend/src/client.gen.ts \
  openapi -out ./openapi.yaml
```

#### Generate server code

- HTTP handler with REST API router
  - `func NewUserStoreServer(serverImplementation UserStore) http.Handler`
  - HTTP handler for all RPC methods
  - Automatic JSON request/response body (un)marshaling
  - Incoming requests call your RPC methods implementation (server logic)
- Sentinel errors that render HTTP codes

```
webrpc-gen -schema=./webrpc.json -target=golang@v0.7.0 -Server -out server/server.gen.go
```

### 3. Mount and serve the API

```go
package main

func main() {
	api := &rpc.API{} // implements interface{}

  handler := rpc.NewPetStoreServer(api)
	http.ListenAndServe(":8080", handler)
}
```

### 4. Implement the `interface{}` (server business logic)

```go
// rpc/user.go
package rpc

func (s *RPC) GetUser(ctx context.Context, uid string) (user *User, err error) {
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

### Enjoy!

..and let us know what you think in [discussions](https://github.com/golang-cz/gospeak/discussions).

# Authors
- [golang.cz](https://golang.cz)
- [VojtechVitek](https://github.com/VojtechVitek)

# License

[MIT license](./LICENSE)
