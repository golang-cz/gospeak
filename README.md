# Golang interface as your schema for client/server communication

The easiest way to communicate to your Golang server over HTTP.

## Example

### 1. Define schema as Golang interface

```go
// rpc/api.go
package rpc

type User struct {
    Uid string
    Name string
}

type ExampleAPI interface{
    GetUser(ctx context.Context, uid string) (user *User, err error)
    ListUsers(ctx context.Context) (users []*User, err error)
    CreateUser(ctx context.Context, userReq *User) (user *User, err error)
    UpdateUser(ctx context.Context, userReq *User) (user *User, err error)
    DeleteUser(ctx context.Context, userReq *User) (err error)
}
```

### 2. Generate webrpc schema from the interface

You can pass a single `.go` file or a folder (Go package) as the schema.

```sh
go2webrpc -schema=./rpc -out webrpc.json
```

### 3. Generate server stub code

Generate server code including:

- HTTP handler for the generated `/rpc/*` REST API routes
  - `func NewExampleAPIServer(serverImplementation ExampleAPI) http.Handler`
  - Automatically (un)marshals JSON request/response body into Go variables
  - Calls your implementation of the method
- Errors that render HTTP codes

```
webrpc-gen -schema=./webrpc.json -target=golang@v0.7.0 -Server -out rpc/server.gen.go
```

### 4. Implement interface methods (server code)

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

### 5. Serve your RPC methods over HTTP

```go
package main

func main() {
    server := &rpc.RPC{
        // DB: databaseModels,
    }

    http.ListenAndServe(":8080", rpc.NewExampleAPIServer(server))
}
```

### 6. Generate clients

Golang client:
```
webrpc-gen -schema=./webrpc.json -target=golang@v0.7.0 -Client -out pkg/example/apiClient.gen.go
```

TypeScript client:
```
webrpc-gen -schema=./webrpc.json -target=typescript@v0.7.0 -Client -out ../frontend/src/exampleApi.gen.ts
```

### 6. Generate documentation

OpenAPI 3.x (Swagger) documentation:
```
webrpc-gen -schema=./webrpc.json -target=openapi@v0.7.0 -out ./openapi.yaml
```

# Authors
- [golang.cz](https://golang.cz)
- [VojtechVitek](https://github.com/VojtechVitek)

# License

[MIT license](./LICENSE)
