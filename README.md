# GoSpeak - Go `interface{}` as your API

What if Go `interface{}` was your schema for service-to-service communication? What if you could generate REST API server code, documentation and strongly typed clients in Go/TypesScript/JavaScript in seconds? What if you could use Go channels over network easily?

Introducing **GoSpeak**, a lightweight JSON alternative to gRPC and Twirp, where Go `interface{}` is your protobuf schema. GoSpeak is built on top of [webrpc](https://github.com/webrpc/webrpc) JSON protocol & code-generation suite.

## Example

1. Define your API schema with Go `interface{}`
2. Install [gospeak](./releases) and [webrpc-gen](https://github.com/webrpc/webrpc/releases)
3. Generate `webrpc.json` schema from the `interface{}`
4. Generate REST API server handlers
5. Implement `interface{}` (server business logic)
6. Serve the REST API
7. Generate strongly typed clients in Go/TypeScript/JavaScript
8. Generate OpenAPI 3.x (Swagger) documentation
9. Enjoy!

### 2. Define your API schema with Go `interface{}`

```go
package schema

type UserStore interface {
	UpsertUser(ctx context.Context, user *User) (*User,  error)
	GetUser(ctx context.Context, ID int64) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	DeleteUser(ctx context.Context, ID int64) error
}

type User struct {
    ID int64
    UID string
    Name string
}
```

### 2. Install gospeak and webrpc-gen

See [gospeak](./releases) and [webrpc-gen](https://github.com/webrpc/webrpc/releases) releases.

### 3. Generate webrpc.json schema from the `interface{}`

You can pass a single `.go` file or a folder (Go package) as the schema.

```sh
gospeak -schema=./rpc -out webrpc.json
```

### 4. Generate REST API server handlers

Generate server code including:

- REST API router
  - `func NewUserStoreServer(serverImplementation UserStore) http.Handler`
  - HTTP handler for all RPC methods
  - Automatic JSON request/response body (un)marshaling
  - Incoming requests call your server implementation
- Sentinel errors that render HTTP codes

```
webrpc-gen -schema=./webrpc.json -target=golang@v0.7.0 -Server -out server/server.gen.go
```

### 5. Implement `interface{}` (server business logic)

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

### 6. Serve the REST API

```go
package main

func main() {
   	rpc := &server.RPC{
		UserStore: map[int64]*server.User{},
        // Data models, DB connection etc.
	}

	apiServer := server.NewUserStoreServer(rpc)
	http.ListenAndServe(":8080", apiServer)
}
```

### 7. Generate API clients

Golang client:
```
webrpc-gen -schema=./webrpc.json -target=golang@v0.7.0 -Client -out pkg/example/apiClient.gen.go
```

TypeScript client:
```
webrpc-gen -schema=./webrpc.json -target=typescript@v0.7.0 -Client -out ../frontend/src/exampleApi.gen.ts
```

### 8. Generate API documentation

OpenAPI 3.x (Swagger) documentation:
```
webrpc-gen -schema=./webrpc.json -target=openapi@v0.7.0 -out ./openapi.yaml
```

### 9. Enjoy!

..and let us know what you think in [discussions](https://github.com/golang-cz/gospeak/discussions).

# Authors
- [golang.cz](https://golang.cz)
- [VojtechVitek](https://github.com/VojtechVitek)

# License

[MIT license](./LICENSE)
