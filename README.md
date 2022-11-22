# GoSpeak - Go `interface{}` as your API

What if Go `interface{}` was your schema for service-to-service communication? What if you could generate REST API server code, documentation and strongly typed clients in Go/TypesScript/JavaScript in seconds? What if you could use Go channels over network easily?

Introducing **GoSpeak**, a lightweight JSON alternative to gRPC and Twirp, where Go `interface{}` is your protobuf schema. GoSpeak uses [webrpc](https://github.com/webrpc/webrpc) protocol and code-generation suite behind the scenes.

## Example

1. Define Go `interface{}` API
2. Generate REST API server (HTTP handlers with JSON)
3. Implement the `interface{}` (server business logic)
4. Serve the REST API
5. Generate strongly typed clients in Go/TypeScript/JavaScript
6. Generate OpenAPI 3.x (Swagger) documentation
7. Enjoy!

### 2. Define Go `interface{}` API schema

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

### 2. Generate webrpc.json schema from the `interface{}`

You can pass a single `.go` file or a folder (Go package) as the schema.

```sh
go2webrpc -schema=./rpc -out webrpc.json
```

### 3. Generate server stub code

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

### 4. Implement the interface (server business logic)

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

### 5. Serve the REST API

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

### 6. Generate API clients

Golang client:
```
webrpc-gen -schema=./webrpc.json -target=golang@v0.7.0 -Client -out pkg/example/apiClient.gen.go
```

TypeScript client:
```
webrpc-gen -schema=./webrpc.json -target=typescript@v0.7.0 -Client -out ../frontend/src/exampleApi.gen.ts
```

### 6. Generate API documentation

OpenAPI 3.x (Swagger) documentation:
```
webrpc-gen -schema=./webrpc.json -target=openapi@v0.7.0 -out ./openapi.yaml
```

# Authors
- [golang.cz](https://golang.cz)
- [VojtechVitek](https://github.com/VojtechVitek)

# License

[MIT license](./LICENSE)
