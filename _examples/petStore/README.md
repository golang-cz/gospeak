# GoSpeak userStore example

1. Read/update [schema/api.go](./schema/api.go)
2. Generate client/server code
    ```
    $ go generate -x ./...

    gospeak -schema=./api.go -interface=ExampleAPI -out=./webrpc.json
    webrpc-gen -schema=./webrpc.json -target=golang -server -pkg=rpc -out=../server.gen.go
    webrpc-gen -schema=./webrpc.json -target=golang -client -pkg=users -out=../../pkg/users/client.gen.go
    ```
3. Run server
    ```
    $ go run ./
    ```
4. Run client tests
    ```
    $ go test ./...
    ```