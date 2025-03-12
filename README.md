***NOTICE:** Experimental. We welcome early user feedback.*

# GoSpeak <!-- omit in toc -->

GoSpeak is a lightweight code generator that enables you to expose your Go methods as REST API, allowing remote execution from various programming languages, including Go and TypeScript. By generating strongly-typed HTTP clients and OpenAPI documentation, GoSpeak provides a streamlined alternative to [gRPC](https://grpc.io) for web applications, using your Go code as the source of truth and JSON as the data format.

| Feature                 | Description |
|-------------------------|-------------|
| **Remote Execution**    | Invoke Go methods remotely from any language that supports HTTP and JSON. |
| **Code Generation**     | Generate strongly-typed HTTP clients in Go, TypeScript, and JavaScript. |
| **OpenAPI Documentation** | Automatically generate OpenAPI 3.x (Swagger) documentation. |
| **Framework Compatibility** | Works with `net/http`, `chi`, `gin`, and `echo`. The generated `http.Handler` integrates with existing handlers and middleware. |
| **Web Compatibility**   | Uses standard HTTP/HTTPS with JSON, ensuring seamless integration with browsers, HTTP clients, proxies, caches, and tools like `cURL`. |

## Language support <!-- omit in toc -->

GoSpeak uses [webrpc](https://github.com/webrpc/webrpc) to generate REST API client and server code using Go templates. The API routes and JSON payload are defined per webrpc data format and can be exported to OpenAPI (Swagger) documentation.

| Language | Code Generation | Requirements |
| -------- | --------------- | ------- |
| [Go](https://github.com/webrpc/gen-golang) | Server and Client| Go 1.22+ |
| [TypeScript](https://github.com/webrpc/gen-typescript) | Client |
| [JavaScript](https://github.com/webrpc/gen-javascript) | Client | ES6 |
| [Kotlin](https://github.com/webrpc/gen-kotlin) | Client | coroutines, moshi, ktor |
| [Dart](https://github.com/webrpc/gen-Dar) | Client | Dart 3.1+ |
| [OpenAPI](https://github.com/webrpc/gen-openapi) | Documentation | OpenAPI 3+ (Swagger) |
| [OpenAPI](https://github.com/webrpc/gen-openapi#generate-clientdocs-via-openapi-generator) | Clients | See list of [code generators](https://github.com/webrpc/gen-openapi#generate-clientdocs-via-openapi-generator) |

# Quick example <!-- omit in toc -->

- [1. Define service API](#1-define-service-api)
- [2. Add webrpc Target Directives](#2-add-webrpc-target-directives)
- [3. Generate Code](#3-generate-code)
- [4. Mount the Code-Generated http.Handler](#4-mount-the-code-generated-httphandler)
- [5. Implement the Server Business Logic](#5-implement-the-server-business-logic)
- [6. Use the Generated Client for Service-To-Service Communication](#6-use-the-generated-client-for-service-to-service-communication)
- [7. Use the Generated Client in Go Tests](#7-use-the-generated-client-in-go-tests)


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

## 2. Add webrpc Target Directives

The following directives will generate Go server and client code with webrpc:

```diff
+//go:webrpc golang -server -pkg=server -out=./server/server.gen.go
+//go:webrpc golang -client -pkg=client -out=./client/example.gen.go
 type PetStore interface {
```

The following will add TypeScript client and OpenAPI 3.x (Swagger) documentation:

```diff
 //go:webrpc golang -server -pkg=server -out=./server/server.gen.go
 //go:webrpc golang -client -pkg=client -out=./client/example.gen.go
+//go:webrpc typescript -client -out=./client/exampleClient.gen.ts
+//go:webrpc openapi -out=./docs/exampleApi.gen.yaml -title=PetStoreAPI
 type PetStore interface {
```

## 3. Generate Code

Run [gospeak](https://github.com/golang-cz/gospeak/releases) binary to generate webrpc code:

```bash
$ gospeak ./proto/api.go
            PetStore => ./server/server.gen.go ✓
            PetStore => ./client/client.gen.go ✓
            PetStore => ./docs/videoApi.gen.yaml ✓
            PetStore => ./client/videoDashboardClient.gen.ts ✓
```

Alternatively, add gospeak as your tool dependency and run it via `go generate`:
```diff
+//go:generate github.com/golang-cz/gospeak/cmd/gospeak .
 package proto
```

```bash
$ go get -tool github.com/golang-cz/gospeak
```

```bash
$ go generate
            PetStore => ./server/server.gen.go ✓
            PetStore => ./client/client.gen.go ✓
            PetStore => ./docs/videoApi.gen.yaml ✓
            PetStore => ./client/videoDashboardClient.gen.ts ✓
```

## 4. Mount the Code-Generated http.Handler

```go
// cmd/petstore/main.go
package main

import "./server"

func main() {
	api := &server.Server{} // your implementation

	handler := server.NewPetStoreServer(api)

	http.ListenAndServe(":8080", handler)
}
```

## 5. Implement the Server Business Logic

The code generated `http.Handler`:

- handles incoming REST API requests
- decodes JSON request body into method argument(s)
- calls your method implementation
- sets proper headers and status code
- encodes method return argument(s) into a JSON response body

What's left for you is the method implementation:

```go
// rpc/server.go
package rpc

type Server struct {
	 // Dependencies like DB connections, logger, configurations, etc.
}
```

```go
// rpc/user.go
package rpc

func (s *Server) GetUser(ctx context.Context, uid string) (user *User, err error) {
	user, err := s.DB.GetUser(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows {
			return nil, proto.ErrNotFound.WithCausef("no such user(%q)", uid)
		}
		return nil, proto.ErrUnexpected.WithCausef("fetch user(%q): %w", uid, err)
	}

	return user, nil
}
```

See [source code](./_examples/petStore/server/pets.go)

## 6. Use the Generated Client for Service-To-Service Communication

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

## 7. Use the Generated Client in Go Tests

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

..and please let us know your thoughts in [discussions](https://github.com/golang-cz/gospeak/discussions).

# Authors <!-- omit in toc -->
- [golang.cz](https://golang.cz)
- [VojtechVitek](https://github.com/VojtechVitek)

# License <!-- omit in toc -->

[MIT license](./LICENSE)
